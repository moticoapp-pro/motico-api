# ERP Reactivo - Especificación del Producto

## 📑 Índice Principal

### Especificaciones
- [Especificación del Producto](./erp_spec.md) - Visión, arquitectura, UX

### Documentación Técnica
- [Arquitectura API](./architecture.md) - Patrón hexagonal, módulos
- [Patrones & Convenciones](./patterns.md) - Convenciones del código
- [Base de Datos](./database.md) - Esquema, migraciones
- [Debugging](./debugging.md) - Troubleshooting y soluciones

---

## Visión General

**ERP Reactivo** es un sistema de gestión empresarial donde el inventario es el corazón. A diferencia de ERPs tradicionales que requieren múltiples pasos manuales, nuestro sistema anticipa necesidades y ejecuta cascadas de acciones automáticas. El usuario interactúa como si hablara con un "empleado super poderoso" que entiende el contexto y sugiere el siguiente paso.

## Problema que Resolvemos

**Hoy:** Usuario crea factura → va a inventario → actualiza stock → genera alertas manualmente (7+ pasos)

**Con ERP Reactivo:** Usuario crea factura → sistema ejecuta en cascada (1 click)

## Arquitectura del Producto

### Núcleo: Event-Driven Architecture
```
Acción Usuario → Evento Disparado → Handlers en Cascada
                                   ├─ Actualiza Inventario
                                   ├─ Genera Alertas
                                   ├─ Sugiere Reórdenes
                                   ├─ Recalcula Reportes
                                   └─ Actualiza Dashboard
```

### Eventos Principales
- **Venta:** Reduce stock, genera factura, alerta si bajo
- **Compra:** Aumenta stock, recalcula caja, recomienda distribuir
- **Transferencia:** Mueve entre almacenes, sincroniza reportes
- **Ajuste:** Manual de correcciones, auditable

## UX/UI: Sugerencias Inteligentes Híbridas

### Opción C (Hybrid) - Implementación
1. **Suggestion Hub Lateral** - Siempre visible, agrupa acciones por criticidad
2. **Action Cards Contextuales** - Aparecen post-acción o en momento oportuno
3. **Inline Smart Buttons** - Dentro de tablas/listados para acciones rápidas

### Priorización Visual
- 🔴 **Críticas** (Stock bajo) → Botones grandes, rojo, pulsante
- 🟡 **Normales** (Oportunidades) → Botones medianos, azul
- ℹ️ **Informativas** (Insights) → Links/badges suave

## Flujo Ideal del Usuario

```
Entrada → Ve Sugerencias → Actúa (1 click) → Sistema Cascada → Dashboard Actualiza → Nuevas Sugerencias
```

## Stack Recomendado

- **Backend:** Go (event-driven, goroutines)
- **Frontend:** React/Next.js (real-time updates)
- **Real-time:** WebSocket o SSE
- **Base de Datos:** PostgreSQL (transacciones)
- **Event Bus:** Redis Streams o simple en-memory

## Módulos Fase 1 (MVP)

1. **Gestión de Inventario** (núcleo)
2. **Gestión de Ventas** (facturación + impacto stock)
3. **Alertas Inteligentes** (stock bajo, margen bajo)
4. **Dashboard Reactivo** (visualización en tiempo real)
5. **Recomendaciones** (reorden automática)

## Métricas de Éxito

- ✅ Usuario realiza tarea en 70% menos clics
- ✅ Stock nunca sorprende (alertas trabajan)
- ✅ Sugerencias aceptadas en 60%+ de casos
- ✅ Dashboard actualiza en <500ms
- ✅ Cascadas completan en <2s

## Siguiente Paso

1. Definir eventos específicos y sus handlers
2. Data model para cascadas (tablas de auditoría)
3. Diseño en Figma (componentes + interacciones)
4. Implementar Event Bus core
