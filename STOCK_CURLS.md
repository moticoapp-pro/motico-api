# Stock API - CURL Commands

## 1. GET Stock Total + Desglose por Store
Obtiene el stock total del producto y desglose por cada store.

```bash
curl -X GET "http://localhost:8080/api/v1/products/{productId}/stock" \
  -H "X-Tenant-ID: {tenantId}" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json"
```

**Respuesta:**
```json
{
  "tenant_id": "...",
  "product_id": "...",
  "total_quantity": 100,
  "total_reserved": 5,
  "available_quantity": 95,
  "stores": [
    {
      "store_id": "...",
      "quantity": 10,
      "reserved_quantity": 2,
      "available_quantity": 8
    },
    {
      "store_id": "...",
      "quantity": 90,
      "reserved_quantity": 3,
      "available_quantity": 87
    }
  ],
  "product": { ... }
}
```

## 2. GET Stock por Store Específico
Obtiene el stock de un producto en un store específico usando query parameter.

```bash
curl -X GET "http://localhost:8080/api/v1/products/{productId}/stock?store_id={storeId}" \
  -H "X-Tenant-ID: {tenantId}" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json"
```

**Respuesta:**
```json
{
  "id": "...",
  "tenant_id": "...",
  "product_id": "...",
  "store_id": "...",
  "quantity": 10,
  "reserved_quantity": 2,
  "available_quantity": 8,
  "updated_at": "...",
  "product": { ... }
}
```

## 3. PUT Update Stock (Requiere store_id en body)
Actualiza la cantidad de stock para un producto en un store específico.

```bash
curl -X PUT "http://localhost:8080/api/v1/products/{productId}/stock" \
  -H "X-Tenant-ID: {tenantId}" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": "{storeId}",
    "quantity": 50
  }'
```

**Body:**
```json
{
  "store_id": "uuid-del-store",
  "quantity": 50
}
```

## 4. PATCH Adjust Stock (Requiere store_id en body)
Ajusta (suma o resta) cantidad de stock para un producto en un store específico.

```bash
curl -X PATCH "http://localhost:8080/api/v1/products/{productId}/stock" \
  -H "X-Tenant-ID: {tenantId}" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": "{storeId}",
    "amount": 10
  }'
```

**Body (positivo suma, negativo resta):**
```json
{
  "store_id": "uuid-del-store",
  "amount": 10
}
```

Ejemplo para restar:
```json
{
  "store_id": "uuid-del-store",
  "amount": -5
}
```
