# Figma Handoff Contract — ERP Reactivo Motico

Documento bidireccional: brief para disenar en Figma + reglas para consumir el diseno en codigo.
Referencia: `FRONTEND_DESIGN_GUIDE.md` (UX) | `FRONTEND_ARCHITECTURE.md` (tech stack).

---

## 1. Estructura del Archivo Figma

El archivo Figma debe organizarse en estas pages (en este orden):

```
Page 1: _Tokens            → Variables, colores, tipografia, espaciado
Page 2: _Components        → Libreria de componentes (Atomic Design)
Page 3: Inventory          → Screens del modulo inventario
Page 4: Transfers          → Screens del modulo transferencias
Page 5: Dashboard          → Screens del dashboard
Page 6: Auth               → Login, onboarding
Page 7: Settings           → Configuracion tenant/usuario
Page 8: Mobile             → Versiones responsive de P1-P5
```

Prefijo `_` = pages internas (no son pantallas, son sistema de diseno).

---

## 2. Tokens Figma → CSS Variables

El disenador DEBE crear Figma Variables (no solo estilos) con estos nombres exactos.
El desarrollador mapeara 1:1 a CSS custom properties.

### Colores (crear como Variable Collection: "Colors")

| Figma Variable | CSS Variable | Dark | Light |
|----------------|-------------|------|-------|
| `bg/base` | `--bg-base` | `#0F1117` | `#F5F5F7` |
| `surface/1` | `--surface-1` | `#1A1D26` | `#FFFFFF` |
| `surface/2` | `--surface-2` | `#22263A` | `#F0F0F5` |
| `border/default` | `--border` | `#2E3247` | `#E0E0E5` |
| `status/critical` | `--critical` | `#FF3B30` | `#FF3B30` |
| `status/warning` | `--warning` | `#FF9500` | `#FF9500` |
| `status/success` | `--success` | `#30D158` | `#30D158` |
| `status/info` | `--info` | `#0A84FF` | `#0A84FF` |
| `text/primary` | `--text-primary` | `#F2F2F7` | `#1C1C1E` |
| `text/secondary` | `--text-secondary` | `#8E8EA0` | `#6E6E80` |

### Tipografia (Variable Collection: "Typography")

| Figma Style | CSS | Font | Size |
|-------------|-----|------|------|
| `data/xl` | `font-data text-2xl` | DM Mono | 32px |
| `data/lg` | `font-data text-xl` | DM Mono | 24px |
| `data/base` | `font-data text-base` | DM Mono | 15px |
| `ui/lg` | `font-ui text-lg` | Inter | 18px |
| `ui/base` | `font-ui text-base` | Inter | 15px |
| `ui/sm` | `font-ui text-sm` | Inter | 13px |
| `ui/xs` | `font-ui text-xs` | Inter | 11px |

### Espaciado (Variable Collection: "Spacing")

| Figma Variable | CSS | Value |
|----------------|-----|-------|
| `space/xs` | `gap-1` / `p-1` | 4px |
| `space/sm` | `gap-2` / `p-2` | 8px |
| `space/md` | `gap-3` / `p-3` | 12px |
| `space/lg` | `gap-4` / `p-4` | 16px |
| `space/xl` | `gap-6` / `p-6` | 24px |
| `space/2xl` | `gap-8` / `p-8` | 32px |

---

## 3. Componentes Figma → React Components

Nombrar componentes en Figma con `/` para jerarquia. El desarrollador usara el mismo
nombre como ruta de archivo.

| Figma Component | React Path | Props esperadas |
|-----------------|-----------|-----------------|
| `Atoms/Badge` | `shared/components/atoms/Badge` | `variant: 'critical' \| 'warning' \| 'success' \| 'info'` |
| `Atoms/Input` | `shared/components/atoms/Input` | `label, error, disabled` |
| `Atoms/Toggle` | `shared/components/atoms/Toggle` | `checked, onChange` |
| `Atoms/Spinner` | `shared/components/atoms/Spinner` | `size: 'sm' \| 'md' \| 'lg'` |
| `Molecules/ActionCard` | `shared/components/molecules/ActionCard` | `priority, title, description, onAction` |
| `Molecules/FilterBar` | `shared/components/molecules/FilterBar` | `filters[], onFilter` |
| `Molecules/DataRow` | `shared/components/molecules/DataRow` | `columns[], actions[]` |
| `Organisms/TopBar` | `shared/components/organisms/TopBar` | `tenant, onSearch, alerts[]` |
| `Organisms/NavSide` | `shared/components/organisms/NavSide` | `items[], collapsed` |
| `Organisms/SuggestionHub` | `shared/components/organisms/SuggestionHub` | `suggestions[]` |
| `Organisms/InventoryTable` | `shared/components/organisms/InventoryTable` | `products[], onSelect` |

### Reglas para el disenador
- Cada componente debe tener **variantes** (estados) como Figma Variants, no layers ocultos
- Nombrar variantes: `State=Default`, `State=Hover`, `State=Active`, `State=Disabled`
- Para Badge: `Variant=Critical`, `Variant=Warning`, `Variant=Success`, `Variant=Info`
- Usar **Auto Layout** en todo. El desarrollador implementa con Flexbox/Grid
- **No usar posicionamiento absoluto** excepto: tooltips, dropdowns, toasts

---

## 4. Screens — Requerimientos por Pantalla

Cada screen debe disenar 3 breakpoints como frames separados:

| Breakpoint | Frame width | Nombre frame |
|------------|------------|--------------|
| Mobile | 375px | `Screen/Mobile` |
| Tablet | 1024px | `Screen/Tablet` |
| Desktop | 1440px | `Screen/Desktop` |

### Screens Sprint 1 (obligatorias)

**4.1 Inventory/List**
- Tabla con columnas: SKU, Producto, Stock (DM Mono), Min, Estado (Badge)
- FilterBar arriba: dropdowns Producto, Categoria, Almacen + boton "+ Nuevo"
- Fila hover: mostrar acciones inline [Transferir] [Ajustar]
- Estado vacio: ilustracion + CTA "Agregar primer producto"

**4.2 Inventory/Detail**
- Side panel (no modal, no page nueva). Width: 400px desktop
- Header: nombre + SKU + boton cerrar
- Body: stock actual (numero grande, DM Mono), stock minimo, ultimo movimiento
- Acciones: [Transferir] [Ajustar] [Ver historial]
- Historial: lista de movimientos con +/- color coded

**4.3 SuggestionHub**
- Panel lateral derecho, 256px width
- Cards apiladas verticalmente, ordenadas por prioridad
- 3 variantes de card: Critical (borde rojo), Warning (borde amarillo), Info (borde azul)
- Cada card: icono + titulo + data + boton accion
- Header: "Acciones Sugeridas (N)" con contador

**4.4 Transfer/Flow**
- Paso 1: Seleccionar origen (almacen + producto + cantidad)
- Paso 2: Seleccionar destino (almacen)
- Paso 3: Confirmar (resumen visual)
- Usar stepper horizontal para indicar progreso
- Boton "Cancelar" siempre visible

**4.5 Stock/Adjust**
- Inline en el panel de detalle (no screen separada)
- Input numerico con +/- buttons
- Selector de razon: Conteo fisico, Dano, Devolucion, Otro
- Preview: "Stock actual: 45 → Nuevo: 43"

---

## 5. Reglas de Handoff (Figma → Codigo)

### Para el disenador
1. Entregar usando **Figma Dev Mode** activado en el archivo
2. Exportar iconos como SVG individuales con nombres `icon-{nombre}.svg`
3. Marcar componentes como "Ready for dev" cuando esten aprobados
4. Anotar interacciones con Figma Comments (no prototyping para MVP)
5. Toda imagen/ilustracion exportable a 1x y 2x

### Para el desarrollador
1. Usar **Figma MCP** (ya configurado) para extraer design context automaticamente
2. Flujo: `get_design_context(nodeId, fileKey)` → adaptar a Next.js + Tailwind
3. NO copiar CSS literal de Figma — mapear a tokens Tailwind definidos en `tokens.css`
4. Componentes con variantes Figma = props con union types en React
5. Spacing de Figma = clases Tailwind del spacing table (seccion 2)
6. Priorizar componentes de `shadcn/ui` como base cuando el diseno lo permita

### Checklist de entrega por screen
- [ ] 3 breakpoints diseñados (375, 1024, 1440)
- [ ] Todos los tokens son Figma Variables (no colores hardcoded)
- [ ] Componentes usan Auto Layout
- [ ] Estados interactivos como Variants (hover, active, disabled, error, empty)
- [ ] Iconos exportables como SVG
- [ ] Dev Mode activado y componentes marcados "Ready for dev"

---

## 6. Flujo de Trabajo Figma ↔ Codigo

```
Disenador                          Desarrollador
    │                                    │
    ├── Disena en Figma ──────────────► Revisa en Figma Dev Mode
    ├── Marca "Ready for dev" ────────► get_design_context (MCP)
    │                                    ├── Extrae tokens + estructura
    │                                    ├── Mapea a componentes existentes
    │                                    └── Implementa con Tailwind + React
    │                                    │
    ◄── Review visual (Vercel preview) ──┤
    ├── Feedback en Figma Comments ───► Ajusta codigo
    └── Aprueba ──────────────────────► Merge
```

Iteracion rapida: el desarrollador usa `get_design_context` del MCP de Figma
para obtener screenshot + code reference + tokens en cada iteracion. No esperar
handoff formal — consumir diseno tan pronto como el frame este marcado ready.
