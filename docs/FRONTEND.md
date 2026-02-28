# Arquitectura Frontend — ERP Reactivo Motico

## 1. Punto de Partida

El backend `motico-api` ya expone dominios funcionales via REST (`/api/v1`):
categories, stores, products, stock, transfers, auth. Arquitectura hexagonal en Go,
multi-tenant con JWT. El frontend debe reflejar estos dominios sin duplicar logica.

---

## 2. Enfoque: Schema-Driven Design (SDD)

El contrato API (Swagger/OpenAPI de motico-api) es la fuente de verdad unica.

```
OpenAPI spec (Go/Swaggo) --> Generador --> Types + API Client (TS)
                                      --> Validaciones (Zod desde schema)
                                      --> Mocks para tests
```

**Herramienta:** `openapi-typescript` + `openapi-fetch` generan tipos y cliente HTTP
directamente del spec. Zero tipos manuales. Cambio en backend = regenerar = error en
compile si el frontend no se adapta.

---

## 3. Stack Tecnologico

| Capa | Tecnologia | Razon |
|------|-----------|-------|
| Framework | **Next.js 15 (App Router)** | SSR/RSC, file-based routing, optimizado para dashboards |
| Lenguaje | **TypeScript strict** | Seguridad de tipos alineada con Go typed backend |
| Styling | **Tailwind CSS 4** | Tokens del Design Guide mapeados a CSS variables |
| Estado servidor | **TanStack Query v5** | Cache, revalidacion, optimistic updates |
| Estado local | **Zustand** | Minimo, sin boilerplate, stores por dominio |
| Real-time | **SSE (Server-Sent Events)** | Go nativo (`net/http`), mas simple que WS para push unidireccional |
| Tablas | **TanStack Table v8** | Headless, virtualizado, sorting/filtering nativo |
| Forms | **React Hook Form + Zod** | Validacion derivada del schema OpenAPI |
| Charts | **Recharts** | Ligero, composable, suficiente para dashboard MVP |
| Generacion API | **openapi-typescript + openapi-fetch** | SDD: tipos y cliente del spec |
| Testing | **Vitest + Testing Library + MSW** | Mocks generados del spec OpenAPI |
| Monorepo | **Turborepo** | Si frontend y backend coexisten en `moticoapp/` |

**Complementariedad con Go:** SSE nativo en Go (goroutines por conexion), OpenAPI generado
por Swaggo, JSON contracts estrictos. No se necesita GraphQL ni BFF adicional.

---

## 4. Estructura del Proyecto (Domain-Driven)

```
motico-web/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── (auth)/             # Grupo: login, registro
│   │   ├── (dashboard)/        # Grupo: shell con nav + hub
│   │   │   ├── inventory/      # /inventory → tabla viva
│   │   │   ├── transfers/      # /transfers → flujo 2 pasos
│   │   │   ├── categories/     # /categories → CRUD
│   │   │   ├── stores/         # /stores → gestion almacenes
│   │   │   └── dashboard/      # /dashboard → metricas
│   │   └── layout.tsx          # Shell: topbar + nav + hub
│   │
│   ├── domains/                # Logica por dominio (mirrors backend)
│   │   ├── inventory/
│   │   │   ├── api.ts          # Queries/mutations (TanStack Query)
│   │   │   ├── store.ts        # Estado local (Zustand)
│   │   │   ├── types.ts        # Re-export de tipos generados
│   │   │   └── hooks.ts        # useInventoryTable, useStockAlerts
│   │   ├── transfers/
│   │   ├── categories/
│   │   ├── stores/
│   │   └── suggestions/        # Dominio del Suggestion Hub
│   │       ├── engine.ts       # Logica de priorizacion
│   │       └── hooks.ts        # useSuggestions
│   │
│   ├── shared/
│   │   ├── api/
│   │   │   ├── client.ts       # openapi-fetch configurado con JWT + tenant
│   │   │   └── generated/      # Tipos auto-generados (NO editar)
│   │   ├── components/         # Atomic Design
│   │   │   ├── atoms/          # Badge, Tag, Input, Toggle, Spinner
│   │   │   ├── molecules/      # DataRow, ActionCard, FilterBar
│   │   │   └── organisms/      # InventoryTable, SuggestionHub, TopBar
│   │   ├── hooks/              # useAuth, useTenant, useSSE
│   │   └── lib/                # Utilidades puras
│   │
│   └── styles/
│       └── tokens.css          # Variables del Design Guide
```

---

## 5. Mapeo Design Guide → Implementacion

### Tokens CSS (desde FRONTEND_DESIGN_GUIDE.md)
```css
:root {
  --bg-base: #0F1117;  --surface-1: #1A1D26;  --surface-2: #22263A;
  --border: #2E3247;   --critical: #FF3B30;    --warning: #FF9500;
  --success: #30D158;  --info: #0A84FF;
  --text-primary: #F2F2F7;  --text-secondary: #8E8EA0;
  --font-data: 'DM Mono', monospace;  --font-ui: 'Inter', sans-serif;
}
.light { --bg-base: #F5F5F7; --surface-1: #FFFFFF; /* ... */ }
```

### Shell Layout (3 columnas en desktop)
```
TopBar:          sticky top, z-50, tenant selector + search + alerts
NavSide:         w-[200px] colapsable a w-[48px], 5 items MVP
Content:         flex-1, scroll interno
SuggestionHub:   w-[256px] fixed right, bottom-sheet en mobile
```

### Responsive segun Guide
- `< 768px`: stack vertical, bottom nav, hub como bottom sheet
- `768-1280px`: 2 cols, nav colapsado, hub overlay
- `> 1280px`: 3 cols completas

---

## 6. Flujo de Datos Real-Time

```
[Go Backend]                    [Frontend]
   │                               │
   ├── POST /transfers ──────────► TanStack Query mutate
   │                               ├── Optimistic update (tabla)
   │                               └── Invalidate queries relacionadas
   │
   ├── SSE /events ──────────────► useSSE hook
   │   (stock.updated,            ├── Actualiza cache TanStack Query
   │    alert.created,            ├── Suggestion Hub re-render
   │    transfer.completed)       └── Animacion numero stock
```

**Go side:** Endpoint SSE con goroutine por conexion. Filtra eventos por tenant_id.
**Frontend:** Hook `useSSE` reconecta automaticamente, parsea eventos tipados.

---

## 7. Multi-Tenancy en Frontend

- Header `X-Tenant-ID` inyectado en cada request via `client.ts`
- Zustand store `useTenantStore` persiste tenant activo en `localStorage`
- Cambio de tenant: invalida todo el cache de TanStack Query
- Subdominio (`tenant.motico.app`) o query param para deep-linking

---

## 8. Suggestion Hub — Motor de Priorizacion

```typescript
// domains/suggestions/engine.ts
type Priority = 'critical' | 'warning' | 'info';
type Suggestion = { id: string; priority: Priority; domain: string; action: () => void };

// Reglas evaluadas contra estado actual del inventario
// Alimentado por SSE events + polling cada 30s como fallback
// Ordenado: critical > warning > info, luego por timestamp
```

El hub no es un modulo backend dedicado — es una **proyeccion del estado de multiples
dominios** calculada en frontend. Esto evita acoplar logica de UX al backend.

---

## 9. Plan de Fases

### Fase 1 — Fundacion (Sprint 1-2)
- Setup Next.js + SDD pipeline (openapi-typescript)
- Shell layout (TopBar, NavSide, Content area)
- Auth flow (login, JWT storage, refresh)
- Inventory table con estados critico/alerta/ok
- Side panel de detalle producto

### Fase 2 — Interactividad (Sprint 3-4)
- Suggestion Hub con cards priorizadas
- Transferencias: flujo 2 pasos origen→destino→confirmar
- Ajuste de stock inline
- SSE para actualizaciones en vivo
- Animaciones de stock (count-up/down)

### Fase 3 — Dashboard + Mobile (Sprint 5-6)
- Dashboard con metricas real-time (Recharts)
- Busqueda global multi-entidad
- Responsive mobile (bottom nav, bottom sheet hub)
- PWA para operacion en piso

---

## 10. Decisiones Arquitectonicas Clave

| Decision | Eleccion | Alternativa descartada | Razon |
|----------|----------|----------------------|-------|
| Real-time | SSE | WebSocket | Unidireccional suficiente, Go nativo, sin lib extra |
| Estado servidor | TanStack Query | SWR | Devtools, optimistic updates, invalidation granular |
| Estado local | Zustand | Redux/Context | Minimo boilerplate, stores aislados por dominio |
| API client | openapi-fetch | Axios/fetch manual | SDD: tipado end-to-end desde spec |
| Routing | App Router (RSC) | Pages Router | Server components para shell, streaming |
| Styling | Tailwind | CSS Modules/Styled | Tokens del Design Guide como variables, utility-first |
| Hub logic | Frontend | Backend endpoint | Es proyeccion UX, no logica de negocio |

---

## 11. Figma Integration

Ver `docs/FIGMA_HANDOFF_CONTRACT.md` para el contrato completo Figma ↔ Codigo.
Resumen: tokens Figma Variables mapeados 1:1 a CSS variables, componentes nombrados
con Atomic Design (`Atoms/Badge` → `shared/components/atoms/Badge`), Figma MCP
para extraer design context automaticamente durante desarrollo.
