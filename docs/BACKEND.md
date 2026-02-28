# Arquitectura Backend — ERP Reactivo Motico

Evolucionar `motico-api` de CRUD multi-tenant a ERP reactivo event-driven.
Mantener arquitectura hexagonal existente. Alineado con frontend Next.js 15 + SSE + TanStack Query.

---

## 1. Que Agregar (el CRUD existente no se reescribe)

| Capacidad | Prioridad | Para que |
|-----------|-----------|----------|
| Event Bus in-process | P0 | Cascadas reactivas, alimenta SSE |
| SSE endpoint | P0 | Frontend `useSSE` → tabla viva, Suggestion Hub |
| Auth real (bcrypt, refresh) | P0 | Login, usuarios, roles |
| Alerts domain | P1 | Datos crudos para Suggestion Hub |
| Audit log | P1 | Historial en panel de detalle |
| Transaction Manager | P1 | Cascadas atomicas |
| Pagination estandar | P1 | TanStack Table server-side |
| Sales domain | P1 | Facturacion + impacto stock |
| RBAC | P2 | admin / operator / viewer |
| Search global | P2 | Busqueda multi-entidad |

---

## 2. Estructura Objetivo

```
internal/
├── domain/
│   ├── auth/          # Evolucionar: users, bcrypt, refresh tokens
│   ├── category/      # Mantener
│   ├── product/       # Mantener, emitir eventos
│   ├── stock/         # Mantener, emitir eventos
│   ├── store/         # Mantener
│   ├── transfer/      # Mantener, emitir eventos
│   ├── tenant/        # Evolucionar: service layer
│   ├── alert/         # NUEVO
│   └── sale/          # NUEVO
├── eventbus/          # NUEVO: bus in-process
│   ├── bus.go         # Interfaz Bus + inMemoryBus
│   ├── payloads.go    # Tipos de payload por evento
│   └── handlers.go    # Handlers de cascada
├── repository/
│   ├── tx.go          # NUEVO: Transaction manager
│   ├── alert.go       # NUEVO
│   ├── sale.go        # NUEVO
│   ├── audit.go       # NUEVO
│   └── ...existentes
└── rest/
    ├── sse/           # NUEVO: SSE handler
    ├── alert/         # NUEVO
    ├── sale/          # NUEVO
    ├── search/        # NUEVO
    └── ...existentes

pkg/
└── pagination/        # NUEVO
```

---

## 3. Event Bus

Bus in-memory con goroutines. Sin Redis/NATS para MVP (migrable despues).

```go
// internal/eventbus/bus.go
type Event struct {
    Type      string    // "stock.updated", "transfer.completed"
    TenantID  string
    Payload   any
    Timestamp time.Time
}

type Bus interface {
    Publish(ctx context.Context, event Event)
    Subscribe(eventType string, handler Handler)
}
```

Servicios existentes emiten eventos al final de cada mutacion (cambio minimo):

```go
func (s *Service) Adjust(ctx context.Context, ...) error {
    // ...logica existente...
    s.bus.Publish(ctx, eventbus.Event{Type: "stock.updated", TenantID: tid, Payload: ...})
    return nil
}
```

### Eventos MVP

| Evento | Origen | Cascada |
|--------|--------|---------|
| `stock.updated` | Stock service | Evaluar alertas → Push SSE |
| `transfer.completed` | Transfer service | Actualizar stocks → Push SSE → Alertas |
| `transfer.cancelled` | Transfer service | Liberar reserva → Push SSE |
| `sale.created` | Sale service | Reducir stock → Alertas → Push SSE |
| `alert.created` | Alert handler | Push SSE al Suggestion Hub |
| `product.created` | Product service | Init stock en 0 → Push SSE |

**Regla**: eventos se emiten DESPUES del commit de la transaccion.

---

## 4. SSE (Real-Time)

```
GET /api/v1/events/stream
Headers: Authorization + X-Tenant-ID
Response: text/event-stream
```

Hub mantiene `map[tenantID][]chan Event`. Goroutine por conexion (Go nativo, sin lib).
Event Bus → SSE Hub via subscribers registrados en `main.go`.

Formato:
```
event: stock.updated
data: {"product_id":"...","quantity":43,"previous":45}
```

Limites: 100 conexiones/tenant, heartbeat 30s, buffer 64 eventos/cliente.

---

## 5. Auth Real

**Tablas**: `users` (tenant_id, email unique/tenant, password_hash bcrypt, role, is_active) + `refresh_tokens` (user_id, token_hash, expires_at).

**JWT Claims**: UserID, TenantID, Email, Role. TenantID viene del JWT (no header) para evitar spoofing.

### Endpoints

```
POST /api/v1/auth/login    → { access_token, refresh_token, user }
POST /api/v1/auth/refresh  → { access_token }
POST /api/v1/auth/logout   → revoca refresh token
GET/POST/PATCH /api/v1/users  → CRUD usuarios (admin only)
```

---

## 6. Alertas (datos para Suggestion Hub)

Backend provee datos crudos. Frontend calcula priorizacion y UX del hub.

```go
type Alert struct {
    ID, TenantID, EntityID uuid.UUID
    Type       AlertType    // stock_low, stock_critical, margin_low
    Priority   Priority     // critical, warning, info
    EntityType string       // "product", "store"
    Title      string
    Data       map[string]any
    Status     AlertStatus  // active, dismissed, resolved
}
```

Auto-generadas via event handler: stock.updated → evaluar vs minimum → crear/resolver alerta.

```
GET   /api/v1/alerts             → alertas activas (filtrable)
GET   /api/v1/alerts/summary     → { critical: 2, warning: 5, info: 3 }
PATCH /api/v1/alerts/:id/dismiss
```

---

## 7. Transaction Manager

```go
// internal/repository/tx.go
func (tm *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
    tx, _ := tm.pool.Begin(ctx)
    defer tx.Rollback(ctx)
    if err := fn(ctx, tx); err != nil { return err }
    return tx.Commit(ctx)
}
```

Cascadas multi-tabla dentro de tx. Eventos emitidos despues del commit.

---

## 8. Audit Log

Tabla `audit_log`: tenant_id, entity_type, entity_id, action, changes (JSONB), user_id, created_at.
Indice compuesto en (tenant_id, entity_type, entity_id, created_at DESC).
Registrado via event handlers (async). Endpoint: `GET /api/v1/audit?entity_type=stock&entity_id=<uuid>`

---

## 9. Ventas

```go
type Sale struct {
    ID, TenantID, StoreID, CreatedBy uuid.UUID
    Items  []SaleItem
    Total  float64
    Status SaleStatus // draft, completed, cancelled
}
```

Cascada: validar stock → reservar (tx) → crear sale → confirmar → reducir stock → eventos.

```
GET/POST /api/v1/sales
GET      /api/v1/sales/:id
PATCH    /api/v1/sales/:id/complete
PATCH    /api/v1/sales/:id/cancel
```

---

## 10. Paginacion / Dashboard / Search

- **Paginacion**: `?page=1&per_page=20&sort=name&order=asc` → `{ data, pagination: {page, per_page, total, total_pages} }`
- **Dashboard**: `GET /api/v1/dashboard/metrics` → inventory counts, transfers, sales, alerts summary
- **Search**: `GET /api/v1/search?q=camiseta` → ILIKE + tsvector (sin Elasticsearch para MVP)

---

## 11. DI en main.go

```go
bus := eventbus.New(logger)
stockService := stock.NewService(stockRepo, bus, logger)
transferService := transfer.NewService(transferRepo, stockService, storeRepo, bus, cfg, logger)

// Handlers de cascada
bus.Subscribe("stock.updated", eventbus.NewStockAlertHandler(alertService, stockRepo))
bus.Subscribe("transfer.completed", eventbus.NewTransferCompletedHandler(stockService))
// SSE bridge
bus.Subscribe("stock.updated", func(ctx context.Context, e eventbus.Event) error {
    sseHub.Broadcast(e.TenantID, e); return nil
})
```

---

## 12. Fases de Implementacion

**Fase 0 — Fundacion**: Tx manager, pagination pkg, logger middleware (TODO), tests unitarios, swagger regen.

**Fase 1 — Auth + Events (Sprint 1)**: Migracion users/refresh_tokens, auth service bcrypt, event bus, emision de eventos en stock/transfer, SSE endpoint.

**Fase 2 — Alertas + Ventas (Sprint 2)**: Migracion alerts/audit_log/sales, alert service + auto-generacion, sale service con cascada, audit log via handlers.

**Fase 3 — Polish (Sprint 3)**: Dashboard metrics, search global, RBAC middleware, rate limiting, swagger completo.

---

## 13. Decisiones Clave

| Decision | Eleccion | Razon |
|----------|----------|-------|
| Event bus | In-memory (goroutines) | Zero infra, migrable a NATS/Redis despues |
| Real-time | SSE nativo Go | Unidireccional suficiente, sin lib extra |
| Transactions | pgx.Tx manual | Consistente con stack actual |
| Auth | JWT + refresh + bcrypt | Stateless, multi-tenant friendly |
| Alertas | Backend datos, frontend UX | Desacoplado |
| Audit | Event handler async | Captura cambios de negocio, no HTTP |
| Pagination | Page-based | TanStack Table compatible |
| Search | ILIKE + tsvector | Sin infra extra para MVP |

---

## 14. Dependencias Nuevas

```
golang.org/x/crypto    // bcrypt
github.com/rs/cors     // CORS
```

Todo lo demas cubierto por stack actual (chi, pgx, zap, jwt-go, validator, uuid).

---

## 15. Produccion

- **CORS**: `app.motico.app` + `localhost:3000`
- **Health**: `GET /health`
- **Shutdown**: stop HTTP → close SSE → drain bus → close DB
- **Swagger**: `swag init -g cmd/api/main.go` (obligatorio por SDD frontend)
