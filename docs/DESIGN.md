# Frontend Design Guide — ERP Reactivo

## Visión de Diseño

**Concepto:** "Control Room" — sientes que operas un sistema vivo, no un formulario.
**Tono:** Industrial-preciso. Datos densos pero legibles. Zero decoración inútil.
**Diferenciador:** El inventario respira. Números cambian en vivo, colores alertan sin que el usuario los busque.

---

## Principios

| Principio | Aplicación |
|-----------|------------|
| **Un vistazo, una decisión** | El estado crítico es inmediatamente visible sin leer |
| **Acción mínima** | Máximo 2 clics para cualquier operación core |
| **Feedback instantáneo** | Cada cambio produce respuesta visual en <200ms |
| **Tolerante al error** | Confirmaciones solo para destructivas. El resto: deshacer |
| **No-digital friendly** | Iconos + texto siempre. Nunca icono solo |

---

## Sistema de Color

```
Background base:   #0F1117  (dark, neutro cálido)
Surface 1:         #1A1D26
Surface 2:         #22263A
Borde:             #2E3247

Crítico:           #FF3B30  (pulsante si urgente)
Advertencia:       #FF9500
Éxito / Positivo:  #30D158
Info / Neutral:    #0A84FF
Texto primario:    #F2F2F7
Texto secundario:  #8E8EA0
```

**Light mode (tablet/escritorio):** misma paleta con fondos invertidos `#F5F5F7` / superficies `#FFFFFF`. El sistema detecta OS preference.

---

## Tipografía

```
Display / Números grandes:  "DM Mono" — legibilidad de datos, carácter técnico
UI / Etiquetas / Botones:   "Inter" — neutralidad funcional
Tamaños clave:
  xs: 11px  sm: 13px  base: 15px  lg: 18px  xl: 24px  2xl: 32px
```

---

## Layout: Shell Permanente

```
┌────────────────────────────────────────────────────┐
│  TOPBAR: Tenant selector | Búsqueda global | Alerts │
├──────┬──────────────────────────────┬───────────────┤
│      │                              │               │
│ NAV  │      CONTENT AREA            │  SUGGESTION   │
│ SIDE │      (cambia por módulo)     │  HUB (256px)  │
│      │                              │               │
│      │                              │               │
└──────┴──────────────────────────────┴───────────────┘
```

- **Nav lateral:** colapsable a iconos (48px) en móvil. 5 items max en MVP.
- **Suggestion Hub:** fijo derecho en desktop. En móvil: bottom sheet deslizable con badge contador.
- **Content Area:** viewport completo, scroll interno por sección.

---

## Módulo Core: Inventario

### Vista principal — Tabla Viva

```
┌─────────────────────────────────────────────────────┐
│ [Producto ▼] [Categoría ▼] [Almacén ▼]  [+ Nuevo]  │
├──────┬─────────────────┬──────┬──────┬──────────────┤
│ SKU  │ Producto        │ Stock│ Min  │ Estado       │
├──────┼─────────────────┼──────┼──────┼──────────────┤
│ A001 │ Camiseta Blanca │  3   │  10  │ ⬤ CRÍTICO   │ ← rojo pulsante
│ A002 │ Jean Azul       │  45  │  5   │ ⬤ OK        │
│ A003 │ Zapato Negro    │  8   │  10  │ ⬤ ALERTA    │ ← amarillo
└──────┴─────────────────┴──────┴──────┴──────────────┘
```

- Stock se actualiza en vivo (WebSocket/SSE) con micro-animación del número
- Click en fila → panel lateral (no modal) con detalle + historial + acciones
- Inline actions: `[Transferir] [Ajustar] [Reordenar]` visibles en hover de fila
- Búsqueda global filtra en tiempo real (debounce 200ms)

### Panel de Detalle (side panel, no modal)

```
┌─────────────────────────────┐
│ ← Jean Azul / SKU A002      │
│ ─────────────────────────── │
│ Stock actual:   45 unidades  │
│ Stock mínimo:   5            │
│ Último mov:     Hace 2h      │
│                              │
│ [Transferir]  [Ajustar]      │
│ [Ver historial]              │
│                              │
│ Historial (últimos 10)       │
│ + 20u  Compra  28/02 14:30  │
│ - 5u   Venta   28/02 10:15  │
└─────────────────────────────┘
```

---

## Suggestion Hub — Anatomía

```
┌─────────────────────────┐
│ ACCIONES SUGERIDAS  (3) │
├─────────────────────────┤
│ 🔴 CRÍTICO              │
│ Reponer Camiseta Blanca │
│ Stock: 3 / Mín: 10      │
│ [Crear Orden]           │
├─────────────────────────┤
│ 🟡 ATENCIÓN             │
│ 3 productos bajo mínimo │
│ [Ver todos]             │
├─────────────────────────┤
│ ℹ️ INSIGHT              │
│ Jean Azul: +30% ventas  │
│ esta semana             │
└─────────────────────────┘
```

- Cards ordenadas por criticidad, auto-actualizadas
- Al ejecutar acción desde hub: animación de "dismiss" + nueva sugerencia entra

---

## Flujos Figma — Screens a diseñar (prioridad)

### Prioridad 1 (Sprint 1)
1. **Inventario — Lista** con estados crítico/alerta/ok
2. **Inventario — Panel detalle** (side panel)
3. **Suggestion Hub** con sus 3 variantes de card
4. **Transferencia** — flujo 2 pasos: origen→destino→confirmar
5. **Ajuste de stock** — inline en panel detalle

### Prioridad 2 (Sprint 2)
6. **Dashboard** — métricas en tiempo real
7. **Crear producto** — form simple, campos mínimos
8. **Alertas** — lista histórica
9. **Búsqueda global** — resultados multi-entidad

### Prioridad 3
10. **Onboarding** — setup tenant nuevo
11. **Settings** — tienda, usuarios, integraciones
12. **Mobile** — versiones adaptadas de P1

---

## Responsive / Multiplataforma

| Breakpoint | Layout | Nav | Hub |
|------------|--------|-----|-----|
| < 768px (móvil) | Stack vertical | Bottom nav 5 items | Bottom sheet |
| 768–1280px (tablet) | 2 columnas | Side nav colapsado | Overlay on demand |
| > 1280px (desktop) | 3 columnas | Side nav expandido | Fijo derecho |

**Mobile first para operación en piso** (ej: tomar inventario físico con celular).

---

## Multi-tenant

- **Selector de tenant** en topbar: dropdown con logo + nombre del negocio
- Cambio de tenant → recarga solo datos, mantiene posición de navegación
- Color accent configurable por tenant (brand identity)
- Subdomain o query param: `tenant.motico.app` o `app.motico.app?t=tenant-id`

---

## Componentes Figma — Kit mínimo

```
Atoms:      Badge (status), Tag, Avatar, Spinner, Toggle, Input
Molecules:  DataRow, ActionCard (hub), FilterBar, SidePanel header
Organisms:  InventoryTable, SuggestionHub, TopBar, NavSide
Templates:  Inventory, Dashboard, DetailPanel
```

---

## Micro-interacciones clave

| Trigger | Respuesta |
|---------|-----------|
| Stock baja de mínimo | Badge rojo aparece + hub card nueva (slide-in) |
| Acción ejecutada | Número de stock se anima (count-up/down) |
| Transferencia exitosa | Toast verde 3s + fila actualiza in-place |
| Error | Shake input + toast rojo con mensaje claro |
| Carga inicial | Skeleton screens, nunca spinner bloqueante |

---

## Accesibilidad mínima

- Contraste mínimo AA (4.5:1) en todo texto
- Tab navigation completa en flujos core
- Aria-labels en todos los iconos standalone
- Focus rings visibles con color de acento del tenant
