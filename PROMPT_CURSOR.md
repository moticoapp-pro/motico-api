# Prompt para Cursor: API de Inventarios Multi-Tenant en Go

## ⚠️ Instrucciones Importantes para la Implementación

### Reglas de Desarrollo

1. **Confirmación entre Pasos**:
   - **SIEMPRE** pedir confirmación explícita del usuario antes de pasar al siguiente paso/fase
   - Al completar una fase, mostrar un resumen y preguntar: "¿Procedo con la siguiente fase?"
   - No avanzar automáticamente sin confirmación

2. **Idiomaticidad en Go**:
   - **SIEMPRE** seguir las convenciones idiomáticas de Go
   - Usar nombres de paquetes en minúsculas, una sola palabra
   - Interfaces con nombres descriptivos (Repository, Service, Handler)
   - Pasar `context.Context` como primer parámetro en funciones I/O
   - **Métodos Privados y Públicos**:
   - Métodos públicos (mayúscula inicial): Solo lo que necesita ser exportado del paquete
   - Métodos privados (minúscula inicial): Funciones auxiliares, helpers, validaciones internas
   - Regla: Si no se usa fuera del paquete, debe ser privado
   - Ejemplo: `func (s *Service) Create()` es público, `func (s *Service) validate()` es privado
- **Comentarios de Código**:
   - **NO agregar comentarios de contexto** que repitan lo que el código ya dice
   - Los comentarios deben explicar el "por qué", no el "qué"
   - Documentar funciones públicas con comentarios descriptivos (go doc)
   - Evitar comentarios obvios como `// Incrementa el contador` sobre `counter++`
   - Comentar solo lógica compleja o decisiones de negocio no obvias
- **Configuración Parametrizable**:
   - **SIEMPRE** usar `config/config.json` para valores parametrizables
   - NO hardcodear valores que puedan cambiar (timeouts, límites, URLs, etc.)
   - Valores que deben ir en config:
     - Timeouts y timeouts de conexión
     - Límites de paginación (default, máximo)
     - Tamaños de pool de conexiones
     - URLs de servicios externos
     - Límites de validación (longitud máxima de strings, etc.)
     - Configuración de logging (niveles, formatos)
   - Secrets y credenciales: Variables de entorno (`.env`), nunca en config.json
   - Ejemplo: `max_page_size: 100` en config, no `const MaxPageSize = 100` en código

3. **Gestión de Contexto**:
   - Si el contexto de la conversación se agota o se pierde información:
   - **LEER NUEVAMENTE** el archivo `PROMPT_CURSOR.md` completo antes de continuar
   - Asegurar que todas las decisiones sigan las especificaciones del documento
   - No asumir información que no esté explícitamente documentada

4. **Conexión con Supabase**:
   - Seguir las instrucciones específicas en la sección "Conexión con Supabase"
   - Validar la conexión antes de continuar con otras implementaciones
   - Usar variables de entorno para credenciales (nunca hardcodear)

---

## Contexto del Proyecto

Necesito crear un API REST en Go para el manejo de inventarios con las siguientes características:

### Requisitos Funcionales

1. **Multi-Tenant**: El sistema debe soportar múltiples empresas (tenants) independientes
2. **Entidades Principales**:
   - **Categorías**: Compartidas a nivel de empresa (tenant)
   - **Productos**: Específicos por sucursal/tienda
3. **Modelo de Negocio**:
   - Una empresa (ej: "Family Motorbiker") puede tener múltiples sedes/sucursales
   - Las categorías se comparten entre todas las sucursales de una empresa
   - Los productos son específicos por cada tienda/sucursal
   - El stock es **único por empresa (tenant)**, no por sucursal
   - Ejemplo: Si "Aures" tiene 5 unidades y "Rincón" tiene 8 unidades, el stock total de la empresa es 13

4. **Operaciones CRUD Requeridas**:
   - **Categorías**: Crear, Listar, Eliminar, Modificar
   - **Productos**: Crear, Listar, Eliminar, Modificar cantidades, Traspaso entre sucursales

### Requisitos Técnicos

- **Lenguaje**: Go (Golang)
- **Arquitectura**: Hexagonal (Ports & Adapters)
- **Base de Datos**: Supabase (PostgreSQL)
- **Alcance Inicial**: Solo categorías y productos (scaffolding base)

---

## Arquitectura Hexagonal - Estructura del Proyecto (Orientada al Dominio)

Estructura idiomática de Go siguiendo estándares de la comunidad y arquitectura hexagonal orientada al dominio:

```
motico-api/
├── cmd/
│   └── api/
│       └── main.go                    # Punto de entrada de la aplicación
├── internal/
│   ├── domain/                        # Capa de Dominio (Entidades y Lógica de Negocio)
│   │   ├── category/
│   │   │   ├── entities/             # Entidades y errores del dominio
│   │   │   │   ├── category.go      # Entidad Category
│   │   │   │   └── errors.go        # Errores específicos del dominio (ErrCategoryNotFound, etc)
│   │   │   ├── repository.go         # Interfaz CategoryRepository
│   │   │   └── service.go            # Lógica de negocio de categorías (recibe config)
│   │   ├── product/
│   │   │   ├── entities/             # Entidades y errores del dominio
│   │   │   │   ├── product.go       # Entidad Product
│   │   │   │   └── errors.go        # Errores específicos del dominio
│   │   │   ├── repository.go         # Interfaz ProductRepository
│   │   │   └── service.go            # Lógica de negocio de productos (recibe config)
│   │   ├── store/
│   │   │   ├── entities/             # Entidades y errores del dominio
│   │   │   │   ├── store.go         # Entidad Store
│   │   │   │   └── errors.go        # Errores específicos del dominio
│   │   │   ├── repository.go       # Interfaz StoreRepository
│   │   │   └── service.go            # Lógica de negocio de sucursales (recibe config)
│   │   ├── stock/
│   │   │   ├── entities/             # Entidades y errores del dominio
│   │   │   │   ├── stock.go         # Entidad Stock
│   │   │   │   └── errors.go        # Errores específicos del dominio
│   │   │   ├── repository.go         # Interfaz StockRepository
│   │   │   └── service.go            # Lógica de negocio de stock (recibe config)
│   │   ├── transfer/
│   │   │   ├── entities/             # Entidades y errores del dominio
│   │   │   │   ├── transfer.go      # Entidad Transfer
│   │   │   │   └── errors.go        # Errores específicos del dominio
│   │   │   ├── repository.go         # Interfaz TransferRepository
│   │   │   └── service.go            # Lógica de negocio de traspasos (recibe config)
│   │   └── tenant/
│   │       ├── entities/              # Entidades y errores del dominio
│   │       │   ├── tenant.go         # Entidad Tenant
│   │       │   └── errors.go         # Errores específicos del dominio
│   │       └── repository.go          # Interfaz TenantRepository
│   ├── repository/                    # Implementaciones de Repositorios (Adaptadores de Salida)
│   │   ├── connection.go             # Conexión a Supabase/PostgreSQL
│   │   ├── category.go               # Implementación CategoryRepository (PostgreSQL)
│   │   ├── product.go                # Implementación ProductRepository (PostgreSQL)
│   │   ├── store.go                  # Implementación StoreRepository (PostgreSQL)
│   │   ├── stock.go                  # Implementación StockRepository (PostgreSQL)
│   │   ├── transfer.go               # Implementación TransferRepository (PostgreSQL)
│   │   └── tenant.go                 # Implementación TenantRepository (PostgreSQL)
│   └── rest/                          # Capa REST (Adaptadores de Entrada)
│       ├── category/                  # Handler de categorías
│       │   ├── entities/              # Entidades de request/response para categorías
│       │   │   └── category.go      # Structs: CreateCategoryRequest, UpdateCategoryRequest, CategoryResponse, ListCategoriesResponse
│       │   ├── list.go               # GET /api/v1/categories - Listar todas las categorías del tenant con paginación
│       │   ├── get_by_id.go          # GET /api/v1/categories/{id} - Obtener una categoría específica por ID
│       │   ├── create.go             # POST /api/v1/categories - Crear una nueva categoría
│       │   ├── update.go             # PUT /api/v1/categories/{id} - Actualizar una categoría completa (todos los campos)
│       │   ├── parcial_update.go     # PATCH /api/v1/categories/{id} - Actualizar campos específicos de la categoría (parcial)
│       │   └── remove.go            # DELETE /api/v1/categories/{id} - Eliminar una categoría (validar que no tenga productos)
│       ├── product/                   # Handler de productos
│       │   ├── entities/              # Entidades de request/response para productos
│       │   │   └── product.go       # Structs: CreateProductRequest, UpdateProductRequest, ProductResponse, ListProductsResponse
│       │   ├── list.go               # GET /api/v1/products - Listar productos con filtros (store_id, category_id) y paginación
│       │   ├── get_by_id.go          # GET /api/v1/products/{id} - Obtener un producto específico por ID con información de stock
│       │   ├── create.go             # POST /api/v1/products - Crear un nuevo producto asociado a una sucursal y categoría
│       │   ├── update.go             # PUT /api/v1/products/{id} - Actualizar un producto completo (todos los campos)
│       │   ├── parcial_update.go     # PATCH /api/v1/products/{id} - Actualizar campos específicos del producto (parcial)
│       │   └── remove.go             # DELETE /api/v1/products/{id} - Eliminar un producto (validar que no tenga stock o traspasos)
│       ├── store/                     # Handler de sucursales
│       │   ├── entities/              # Entidades de request/response para sucursales
│       │   │   └── store.go          # Structs: CreateStoreRequest, UpdateStoreRequest, StoreResponse, ListStoresResponse
│       │   ├── list.go                # GET /api/v1/stores - Listar todas las sucursales del tenant con paginación
│       │   ├── get_by_id.go           # GET /api/v1/stores/{id} - Obtener una sucursal específica por ID
│       │   ├── create.go              # POST /api/v1/stores - Crear una nueva sucursal para el tenant
│       │   ├── update.go             # PUT /api/v1/stores/{id} - Actualizar una sucursal completa (todos los campos)
│       │   ├── parcial_update.go     # PATCH /api/v1/stores/{id} - Actualizar campos específicos de la sucursal (parcial)
│       │   └── remove.go             # DELETE /api/v1/stores/{id} - Eliminar una sucursal (validar que no tenga productos)
│       ├── stock/                     # Handler de stock
│       │   ├── entities/              # Entidades de request/response para stock
│       │   │   └── stock.go          # Structs: UpdateStockRequest, AdjustStockRequest, StockResponse
│       │   ├── get_by_id.go          # GET /api/v1/products/{id}/stock - Obtener el stock total del producto a nivel de tenant
│       │   ├── update.go             # PUT /api/v1/products/{id}/stock - Actualizar la cantidad total de stock del producto
│       │   └── adjust.go             # PATCH /api/v1/products/{id}/stock - Ajustar cantidad de stock (incrementar o decrementar)
│       ├── transfer/                  # Handler de traspasos
│       │   ├── entities/              # Entidades de request/response para traspasos
│       │   │   └── transfer.go       # Structs: CreateTransferRequest, UpdateTransferRequest, TransferResponse, ListTransfersResponse
│       │   ├── list.go               # GET /api/v1/transfers - Listar traspasos con filtros (status, store_id) y paginación
│       │   ├── get_by_id.go          # GET /api/v1/transfers/{id} - Obtener un traspaso específico por ID con detalles
│       │   ├── create.go             # POST /api/v1/transfers - Crear un traspaso entre sucursales (validar stock disponible)
│       │   ├── update.go             # PUT /api/v1/transfers/{id} - Actualizar un traspaso (solo si está en estado 'pending')
│       │   ├── complete.go           # PATCH /api/v1/transfers/{id}/complete - Completar un traspaso (cambiar estado a 'completed' y actualizar stock)
│       │   ├── cancel.go             # PATCH /api/v1/transfers/{id}/cancel - Cancelar un traspaso (cambiar estado a 'cancelled' y liberar stock reservado)
│       │   └── remove.go             # DELETE /api/v1/transfers/{id} - Eliminar un traspaso (solo si está en estado 'pending')
│       ├── middleware.go             # Middlewares HTTP compartidos (tenant, auth, logger, recovery)
│       ├── router.go                  # Configuración de rutas
│       └── response.go               # Helpers para respuestas HTTP
├── pkg/                               # Código reutilizable público (puede ser usado por otros proyectos)
│   ├── logger/                        # Logger configurado (wrapper sobre zap/logrus)
│   └── validator/                     # Validadores personalizados
├── migrations/                        # Migraciones de base de datos
│   └── 001_initial_schema.sql
├── config/                            # Configuración de la aplicación
│   ├── config.go                      # Estructura de configuración
│   └── config.json                    # Archivo de configuración JSON
├── config.example.json                 # Ejemplo de configuración
├── .env.example                       # Ejemplo de variables de entorno (para secrets)
├── go.mod
├── go.sum
└── README.md
```

### Principios de la Estructura

1. **Orientación al Dominio**: Cada agregado del dominio tiene su propio paquete con:
   - `entities/`: Entidades y errores del dominio
   - `repository.go`: Interfaz del repositorio
   - `service.go`: Lógica de negocio
2. **Separación de Responsabilidades**:
   - `domain/`: Lógica de negocio pura (sin dependencias externas), incluye errores del dominio
   - `repository/`: Implementaciones de persistencia (sin separación por tecnología)
   - `rest/`: Capa de presentación HTTP (handlers y middlewares juntos para mejor cohesión)
3. **Idiomático en Go**:
   - Uso de `internal/` para código privado de la aplicación
   - Uso de `cmd/` para ejecutables
   - Uso de `pkg/` para código reutilizable público (no es una librería, pero puede ser compartido)
   - Paquetes pequeños y cohesivos
   - Inyección de dependencias mediante constructores que reciben configuraciones
4. **Performance y Escalabilidad**:
   - Repositorio único sin separación por tecnología (fácil cambio de DB sin refactor masivo)
   - Handlers y middlewares en el mismo paquete para mejor performance (menos imports)
   - Configuración inyectada a servicios para mejor testabilidad

---

## Modelo de Datos (Esquema de Base de Datos)

### Tablas Principales

```sql
-- Tenants (Empresas)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Sucursales (Stores/Branches)
CREATE TABLE stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Categorías (Compartidas por tenant)
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Productos (Específicos por sucursal)
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sku VARCHAR(100),
    price DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, store_id, sku)
);

-- Stock (Único por tenant, no por sucursal)
CREATE TABLE stock (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, product_id)
);

-- Traspasos entre sucursales
CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    from_store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    to_store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, completed, cancelled
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CHECK (from_store_id != to_store_id),
    CHECK (status IN ('pending', 'completed', 'cancelled'))
);

-- Índices para optimización
CREATE INDEX idx_products_tenant_store ON products(tenant_id, store_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_stock_tenant_product ON stock(tenant_id, product_id);
CREATE INDEX idx_categories_tenant ON categories(tenant_id);
CREATE INDEX idx_stores_tenant ON stores(tenant_id);
CREATE INDEX idx_transfers_tenant ON transfers(tenant_id);
CREATE INDEX idx_transfers_product ON transfers(product_id);
```

---

## Endpoints REST de la API

### Base Path: `/api/v1`

Todas las rutas requieren:
- Header `X-Tenant-ID` para identificar el tenant
- Header `Authorization: Bearer {token}` para autenticación JWT (excepto endpoint de login)

### Autenticación

```
POST   /api/v1/auth/login              # Iniciar sesión y obtener token JWT
```

**Request:**
```json
{
  "email": "usuario@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### Categorías

```
GET    /api/v1/categories                    # Listar todas las categorías del tenant
GET    /api/v1/categories/{id}               # Obtener categoría por ID
POST   /api/v1/categories                    # Crear nueva categoría
PUT    /api/v1/categories/{id}               # Actualizar categoría completa
PATCH  /api/v1/categories/{id}               # Actualizar categoría parcial
DELETE /api/v1/categories/{id}               # Eliminar categoría
```

### Sucursales

```
GET    /api/v1/stores                        # Listar sucursales del tenant
GET    /api/v1/stores/{id}                   # Obtener sucursal por ID
POST   /api/v1/stores                        # Crear nueva sucursal
PUT    /api/v1/stores/{id}                   # Actualizar sucursal completa
PATCH  /api/v1/stores/{id}                   # Actualizar sucursal parcial
DELETE /api/v1/stores/{id}                   # Eliminar sucursal
```

### Productos

```
GET    /api/v1/products                      # Listar productos (query params: ?store_id=xxx&category_id=xxx)
GET    /api/v1/products/{id}                 # Obtener producto por ID
POST   /api/v1/products                      # Crear nuevo producto
PUT    /api/v1/products/{id}                 # Actualizar producto completo
PATCH  /api/v1/products/{id}                 # Actualizar producto parcial
DELETE /api/v1/products/{id}                 # Eliminar producto
```

### Stock

```
GET    /api/v1/products/{id}/stock           # Obtener stock del producto (total por tenant)
PUT    /api/v1/products/{id}/stock           # Actualizar cantidad de stock
PATCH  /api/v1/products/{id}/stock           # Ajustar cantidad de stock (incremento/decremento)
```

### Traspasos

```
GET    /api/v1/transfers                     # Listar traspasos del tenant (query params: ?status=pending&store_id=xxx)
GET    /api/v1/transfers/{id}                # Obtener traspaso por ID
POST   /api/v1/transfers                     # Crear traspaso entre sucursales (estado: 'pending')
PUT    /api/v1/transfers/{id}                # Actualizar traspaso
PATCH  /api/v1/transfers/{id}/complete       # Completar traspaso (cambia estado a 'completed')
PATCH  /api/v1/transfers/{id}/cancel         # Cancelar traspaso (cambia estado a 'cancelled')
DELETE /api/v1/transfers/{id}                # Eliminar traspaso (solo si está en estado 'pending')
```

**Estados de Traspaso:**
- `pending`: Creado, pendiente de procesamiento
- `completed`: Completado exitosamente
- `cancelled`: Cancelado antes de completarse

### Convenciones REST

- **GET**: Recuperar recursos (listar o obtener por ID)
- **POST**: Crear nuevos recursos
- **PUT**: Actualización completa del recurso
- **PATCH**: Actualización parcial del recurso o acciones específicas
- **DELETE**: Eliminar recursos
- **Query Parameters**: Para filtros y paginación (`?page=1&limit=20&store_id=xxx`)
- **Path Parameters**: Para identificar recursos específicos (`{id}`)
- **Status Codes**:
  - `200 OK`: Operación exitosa
  - `201 Created`: Recurso creado
  - `204 No Content`: Eliminación exitosa
  - `400 Bad Request`: Error de validación
  - `401 Unauthorized`: Token JWT inválido o ausente
  - `403 Forbidden`: Sin permisos para la operación
  - `404 Not Found`: Recurso no encontrado
  - `409 Conflict`: Conflicto de negocio (ej: stock insuficiente)
  - `500 Internal Server Error`: Error del servidor

---

## To-Do de Implementación

### Fase 1: Configuración Inicial y Base de Datos
- [ ] Inicializar proyecto Go con `go mod init`
- [ ] Crear estructura de configuración en `config/config.go`
- [ ] Crear `config/config.json` con estructura de configuración
- [ ] Crear `config.example.json` como plantilla
- [ ] Implementar carga de configuración desde JSON
- [ ] **Configurar Pre-commit (VER SECCIÓN ESPECÍFICA ABAJO)**
- [ ] **Conexión con Supabase (VER SECCIÓN ESPECÍFICA ABAJO)**
- [ ] Crear migraciones SQL para las tablas
- [ ] Configurar logger (inyectable a servicios)
- [ ] Configurar manejador de errores personalizado
- [ ] Implementar autenticación JWT:
  - [ ] Instalar dependencia `github.com/golang-jwt/jwt/v5`
  - [ ] Crear servicio de autenticación para generar y validar tokens JWT
  - [ ] Implementar endpoint de login (`POST /api/v1/auth/login`)
  - [ ] Configurar secret key para JWT en variables de entorno
  - [ ] Implementar middleware de autenticación JWT
  - [ ] Proteger endpoints con middleware de autenticación
- [ ] **PEDIR CONFIRMACIÓN antes de continuar a Fase 2**

---

### 🔌 Conexión con Supabase - Instrucciones Detalladas

**OBJETIVO**: Establecer conexión segura y eficiente con Supabase (PostgreSQL)

#### Paso 1: Obtener Credenciales de Supabase
- [ ] Acceder al dashboard de Supabase: https://app.supabase.com
- [ ] Seleccionar o crear el proyecto
- [ ] Ir a **Settings** → **Database**
- [ ] Copiar la **Connection String** (URI de conexión)
- [ ] Obtener los siguientes valores:
  - `DB_HOST`: Host de la base de datos (ej: `db.xxxxx.supabase.co`)
  - `DB_PORT`: Puerto (por defecto: `5432`)
  - `DB_USER`: Usuario (por defecto: `postgres`)
  - `DB_PASSWORD`: Contraseña del proyecto
  - `DB_NAME`: Nombre de la base de datos (por defecto: `postgres`)
  - `DB_SSLMODE`: Modo SSL (recomendado: `require` para producción)

#### Paso 2: Configurar Variables de Entorno
- [ ] Crear archivo `.env` en la raíz del proyecto (NO commitear)
- [ ] Agregar al `.gitignore`: `.env`
- [ ] Crear `.env.example` con estructura sin valores sensibles:
```env
# Supabase Database Configuration
DB_HOST=db.xxxxx.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=postgres
DB_SSLMODE=require
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_CONN_MAX_LIFETIME=5m
```

#### Paso 3: Instalar Dependencias
- [ ] Instalar `github.com/jackc/pgx/v5` y `github.com/jackc/pgx/v5/pgxpool`
- [ ] Instalar `github.com/joho/godotenv` para cargar variables de entorno
```bash
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/joho/godotenv
```

#### Paso 4: Implementar Conexión en `repository/connection.go`
- [ ] Crear función `NewConnectionPool(ctx context.Context, config *Config) (*pgxpool.Pool, error)`
- [ ] Construir connection string desde variables de entorno:
```go
// Formato: postgres://user:password@host:port/dbname?sslmode=require
connString := fmt.Sprintf(
    "postgres://%s:%s@%s:%s/%s?sslmode=%s",
    config.DB.User,
    config.DB.Password,
    config.DB.Host,
    config.DB.Port,
    config.DB.Name,
    config.DB.SSLMode,
)
```
- [ ] Configurar pool de conexiones:
  - `MaxConns`: Máximo de conexiones (desde config)
  - `MaxIdleConns`: Conexiones inactivas máximas
  - `ConnMaxLifetime`: Tiempo máximo de vida de conexión
- [ ] Validar conexión con `pool.Ping(ctx)`
- [ ] Retornar error si la conexión falla
- [ ] Implementar función `Close()` para cerrar el pool correctamente

#### Paso 5: Estructura de Configuración
- [ ] Agregar estructura `DatabaseConfig` en `config/config.go`:
```go
type DatabaseConfig struct {
    Host            string `json:"host"`
    Port            string `json:"port"`
    User            string `json:"user"`
    Password        string `json:"-"` // No serializar en JSON
    Name            string `json:"name"`
    SSLMode         string `json:"ssl_mode"`
    MaxConnections  int    `json:"max_connections"`
    MaxIdleConns    int    `json:"max_idle_connections"`
    ConnMaxLifetime string `json:"conn_max_lifetime"`
}
```
- [ ] Cargar valores desde `.env` usando `godotenv.Load()`
- [ ] Validar que todas las variables requeridas estén presentes

#### Paso 6: Test de Conexión
- [ ] Crear función de prueba `TestConnection()` en `repository/connection_test.go`
- [ ] Verificar que el pool se crea correctamente
- [ ] Ejecutar query simple: `SELECT 1`
- [ ] Validar que la conexión responde
- [ ] **IMPORTANTE**: Si la conexión falla, NO continuar. Resolver el problema primero.

#### Paso 7: Integración en `main.go`
- [ ] Inicializar pool de conexiones al inicio de la aplicación
- [ ] Pasar el pool a los repositorios mediante constructores
- [ ] Implementar graceful shutdown para cerrar conexiones al terminar
- [ ] Agregar logging de estado de conexión (conectado/desconectado)

#### Checklist de Validación
- [ ] ✅ Pool de conexiones se crea sin errores
- [ ] ✅ Query de prueba (`SELECT 1`) ejecuta correctamente
- [ ] ✅ Variables de entorno cargadas correctamente
- [ ] ✅ `.env` está en `.gitignore`
- [ ] ✅ `.env.example` documentado sin valores sensibles
- [ ] ✅ Conexión usa SSL (sslmode=require)
- [ ] ✅ Pool configurado con límites apropiados
- [ ] ✅ Logging de conexión implementado

#### Errores Comunes a Evitar
- ❌ NO hardcodear credenciales en el código
- ❌ NO usar `database/sql` (usar `pgx` directamente)
- ❌ NO crear múltiples pools de conexión
- ❌ NO olvidar cerrar el pool en shutdown
- ❌ NO usar conexiones sin pool (ineficiente)

#### Siguiente Paso
- [ ] **PEDIR CONFIRMACIÓN** al usuario que la conexión funciona correctamente
- [ ] Solo después de confirmación, proceder con migraciones de base de datos

---

### 🔧 Pre-commit Setup - Instrucciones Detalladas

**OBJETIVO**: Configurar hooks de pre-commit para validar código automáticamente antes de cada commit

#### Paso 1: Instalar Pre-commit
- [ ] Instalar pre-commit: `brew install pre-commit` (macOS) o `pip install pre-commit` (Linux/Windows)
- [ ] Verificar instalación: `pre-commit --version`

#### Paso 2: Crear Archivo `.pre-commit-config.yaml`
- [ ] Crear `.pre-commit-config.yaml` en la raíz del proyecto con:
```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
      - id: check-json
      - id: check-toml
      - id: check-merge-conflict
      - id: check-case-conflict

  - repo: local
    hooks:
      - id: go-fmt
        name: go fmt
        entry: bash -c 'go fmt ./...'
        language: system
        types: [go]
        pass_filenames: false
        always_run: true

      - id: go-vet
        name: go vet
        entry: bash -c 'go vet ./...'
        language: system
        types: [go]
        pass_filenames: false
        always_run: true

      - id: go-test
        name: go test
        entry: bash -c 'go test ./... -short'
        language: system
        types: [go]
        pass_filenames: false
        always_run: true

      - id: golangci-lint
        name: golangci-lint
        entry: bash -c 'golangci-lint run ./...'
        language: system
        types: [go]
        pass_filenames: false
        always_run: true
        require_serial: true
```

#### Paso 3: Instalar Hooks
- [ ] Instalar hooks: `pre-commit install`
- [ ] Verificar instalación: `pre-commit --version`
- [ ] Probar manualmente: `pre-commit run --all-files`

#### Paso 4: Configurar golangci-lint (si no está instalado)
- [ ] Instalar golangci-lint:
  - macOS: `brew install golangci-lint`
  - Linux: `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2`
  - Windows: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- [ ] Crear `.golangci.yml` en la raíz del proyecto:
```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unparam
    - gocritic
    - gosec

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style
  goimports:
    local-prefixes: motico-api

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

#### Paso 5: Validar Configuración
- [ ] Ejecutar pre-commit manualmente: `pre-commit run --all-files`
- [ ] Verificar que todos los hooks pasan
- [ ] Hacer un commit de prueba para verificar que los hooks se ejecutan automáticamente
- [ ] Si algún hook falla, corregir el código antes de poder hacer commit

#### Comandos Útiles
- `pre-commit run --all-files`: Ejecutar todos los hooks manualmente
- `pre-commit run`: Ejecutar hooks solo en archivos staged
- `pre-commit uninstall`: Desinstalar hooks (si es necesario)
- `pre-commit autoupdate`: Actualizar hooks a última versión

#### Checklist de Validación
- [ ] ✅ Pre-commit instalado y funcionando
- [ ] ✅ Hooks se ejecutan automáticamente en cada commit
- [ ] ✅ `go fmt` se ejecuta automáticamente
- [ ] ✅ `go vet` valida código sin errores
- [ ] ✅ `golangci-lint` pasa sin errores
- [ ] ✅ Tests se ejecutan antes de commit
- [ ] ✅ Commits fallan si hay errores de validación

#### Notas Importantes
- Los hooks se ejecutan automáticamente en cada `git commit`
- Si un hook falla, el commit se cancela automáticamente
- Corregir errores antes de intentar commit nuevamente
- Para saltar hooks (NO recomendado): `git commit --no-verify`

### Fase 2: Capa de Dominio (Domain Layer)

**⚠️ ANTES DE CONTINUAR**: Confirmar que la Fase 1 está completa y la conexión con Supabase funciona correctamente.
- [ ] Crear paquete `domain/tenant/`:
  - [ ] Entidad `Tenant` en `entities/tenant.go`
  - [ ] Errores en `entities/errors.go`
  - [ ] Interfaz `TenantRepository`
- [ ] Crear paquete `domain/store/`:
  - [ ] Entidad `Store` en `entities/store.go`
  - [ ] Errores en `entities/errors.go`
  - [ ] Interfaz `StoreRepository`
  - [ ] Servicio `StoreService` con lógica de negocio
- [ ] Crear paquete `domain/category/`:
  - [ ] Entidad `Category` en `entities/category.go`
  - [ ] Errores en `entities/errors.go`
  - [ ] Interfaz `CategoryRepository`
  - [ ] Servicio `CategoryService` con lógica de negocio
- [ ] Crear paquete `domain/product/`:
  - [ ] Entidad `Product` en `entities/product.go`
  - [ ] Errores en `entities/errors.go`
  - [ ] Interfaz `ProductRepository`
  - [ ] Servicio `ProductService` con lógica de negocio
- [ ] Crear paquete `domain/stock/`:
  - [ ] Entidad `Stock` en `entities/stock.go`
  - [ ] Errores en `entities/errors.go`
  - [ ] Interfaz `StockRepository`
  - [ ] Servicio `StockService` con cálculo de stock total por tenant
- [ ] Crear paquete `domain/transfer/`:
  - [ ] Entidad `Transfer` en `entities/transfer.go` con estados (pending, completed, cancelled)
  - [ ] Errores en `entities/errors.go`
  - [ ] Interfaz `TransferRepository`
  - [ ] Servicio `TransferService` con lógica de traspasos (recibe config por constructor)

### Fase 3: Adaptadores de Salida (Repository Layer)

**⚠️ ANTES DE CONTINUAR**: Confirmar que la Fase 2 está completa y todas las entidades del dominio están definidas.
- [ ] Configurar conexión a PostgreSQL/Supabase en `repository/connection.go` (pool de conexiones)
- [ ] Implementar `TenantRepository` en `repository/tenant.go` (sin apellido _repository)
- [ ] Implementar `StoreRepository` en `repository/store.go`
- [ ] Implementar `CategoryRepository` en `repository/category.go`
- [ ] Implementar `ProductRepository` en `repository/product.go`
- [ ] Implementar `StockRepository` en `repository/stock.go`
- [ ] Implementar `TransferRepository` en `repository/transfer.go`
- [ ] Implementar transacciones para operaciones complejas (traspasos)
- [ ] Agregar manejo de errores específicos de base de datos
- [ ] Inyectar configuración a repositorios mediante constructores

### Fase 4: Capa REST (REST Layer)

**⚠️ ANTES DE CONTINUAR**: Confirmar que la Fase 3 está completa y todos los repositorios están implementados y probados.
- [ ] Configurar router HTTP (recomendado: `github.com/go-chi/chi/v5`)
- [ ] Crear entidades de request/response en `rest/{resource}/entities/{resource}.go`:
  - [ ] `rest/category/entities/category.go`: Structs para requests/responses de categorías (CreateCategoryRequest, UpdateCategoryRequest, CategoryResponse, ListCategoriesResponse)
  - [ ] `rest/product/entities/product.go`: Structs para requests/responses de productos (CreateProductRequest, UpdateProductRequest, ProductResponse, ListProductsResponse)
  - [ ] `rest/store/entities/store.go`: Structs para requests/responses de sucursales (CreateStoreRequest, UpdateStoreRequest, StoreResponse, ListStoresResponse)
  - [ ] `rest/stock/entities/stock.go`: Structs para requests/responses de stock (UpdateStockRequest, AdjustStockRequest, StockResponse)
  - [ ] `rest/transfer/entities/transfer.go`: Structs para requests/responses de traspasos (CreateTransferRequest, UpdateTransferRequest, TransferResponse, ListTransfersResponse)
- [ ] Implementar middlewares en `rest/middleware.go`:
  - [ ] `TenantMiddleware`: Extraer `X-Tenant-ID` del header y validar
  - [ ] `AuthMiddleware`: Autenticación JWT (validar token)
  - [ ] `LoggerMiddleware`: Logging de requests
  - [ ] `RecoveryMiddleware`: Manejo de panics
- [ ] Implementar handlers por paquete en `rest/` sin apellidos (el paquete ya indica el recurso):
  - [ ] Paquete `rest/category/`:
    - [ ] `list.go`: GET /api/v1/categories - Listar todas las categorías del tenant con paginación
    - [ ] `get_by_id.go`: GET /api/v1/categories/{id} - Obtener una categoría específica por ID
    - [ ] `create.go`: POST /api/v1/categories - Crear una nueva categoría
    - [ ] `update.go`: PUT /api/v1/categories/{id} - Actualizar una categoría completa
    - [ ] `parcial_update.go`: PATCH /api/v1/categories/{id} - Actualizar campos específicos de la categoría (parcial)
    - [ ] `remove.go`: DELETE /api/v1/categories/{id} - Eliminar una categoría (validar que no tenga productos)
  - [ ] Paquete `rest/store/`:
    - [ ] `list.go`: GET /api/v1/stores - Listar todas las sucursales del tenant
    - [ ] `get_by_id.go`: GET /api/v1/stores/{id} - Obtener una sucursal específica por ID
    - [ ] `create.go`: POST /api/v1/stores - Crear una nueva sucursal
    - [ ] `update.go`: PUT /api/v1/stores/{id} - Actualizar una sucursal completa
    - [ ] `parcial_update.go`: PATCH /api/v1/stores/{id} - Actualizar campos específicos de la sucursal (parcial)
    - [ ] `remove.go`: DELETE /api/v1/stores/{id} - Eliminar una sucursal (validar que no tenga productos)
  - [ ] Paquete `rest/product/`:
    - [ ] `list.go`: GET /api/v1/products - Listar productos con filtros (store_id, category_id) y paginación
    - [ ] `get_by_id.go`: GET /api/v1/products/{id} - Obtener un producto específico por ID con información de stock
    - [ ] `create.go`: POST /api/v1/products - Crear un nuevo producto asociado a una sucursal y categoría
    - [ ] `update.go`: PUT /api/v1/products/{id} - Actualizar un producto completo
    - [ ] `parcial_update.go`: PATCH /api/v1/products/{id} - Actualizar campos específicos del producto (parcial)
    - [ ] `remove.go`: DELETE /api/v1/products/{id} - Eliminar un producto (validar que no tenga stock o traspasos)
  - [ ] Paquete `rest/stock/`:
    - [ ] `get_by_id.go`: GET /api/v1/products/{id}/stock - Obtener el stock total del producto a nivel de tenant
    - [ ] `update.go`: PUT /api/v1/products/{id}/stock - Actualizar la cantidad total de stock del producto
    - [ ] `adjust.go`: PATCH /api/v1/products/{id}/stock - Ajustar cantidad de stock (incrementar o decrementar)
  - [ ] Paquete `rest/transfer/`:
    - [ ] `list.go`: GET /api/v1/transfers - Listar traspasos con filtros (status, store_id) y paginación
    - [ ] `get_by_id.go`: GET /api/v1/transfers/{id} - Obtener un traspaso específico por ID con detalles
    - [ ] `create.go`: POST /api/v1/transfers - Crear un traspaso entre sucursales (validar stock disponible)
    - [ ] `update.go`: PUT /api/v1/transfers/{id} - Actualizar un traspaso (solo si está en estado 'pending')
    - [ ] `complete.go`: PATCH /api/v1/transfers/{id}/complete - Completar un traspaso (cambiar estado y actualizar stock)
    - [ ] `cancel.go`: PATCH /api/v1/transfers/{id}/cancel - Cancelar un traspaso (cambiar estado y liberar stock reservado)
    - [ ] `remove.go`: DELETE /api/v1/transfers/{id} - Eliminar un traspaso (solo si está en estado 'pending')
- [ ] Cada handler recibe servicio y config por constructor
- [ ] Configurar validación de requests (usar `github.com/go-playground/validator/v10`)
- [ ] Implementar helpers de respuesta HTTP en `rest/response.go`
- [ ] Configurar manejo de errores HTTP con códigos de estado apropiados
- [ ] Implementar inyección de dependencias en `cmd/api/main.go` (servicios → handlers)

### Fase 5: Lógica de Negocio Específica

**⚠️ ANTES DE CONTINUAR**: Confirmar que la Fase 4 está completa y todos los endpoints REST están implementados.
- [ ] Implementar validación: categorías compartidas por tenant
- [ ] Implementar validación: productos específicos por sucursal
- [ ] Implementar cálculo de stock total por tenant (suma de todas las sucursales)
- [ ] Implementar lógica de traspaso:
  - [ ] Validar stock disponible al crear traspaso
  - [ ] Reservar stock cuando traspaso está en estado 'pending'
  - [ ] Actualizar stock al completar traspaso
  - [ ] Liberar stock al cancelar traspaso
- [ ] Implementar validaciones de integridad referencial
- [ ] Inyectar configuraciones a todos los servicios mediante constructores

### Fase 6: Testing y Documentación

**⚠️ ANTES DE CONTINUAR**: Confirmar que la Fase 5 está completa y toda la lógica de negocio está implementada.
- [ ] Crear tests unitarios para servicios de dominio
- [ ] Crear tests de integración para repositorios
- [ ] Crear tests de integración para handlers HTTP
- [ ] Configurar Swagger/OpenAPI:
  - [ ] Integrar `github.com/swaggo/swag` o `github.com/getkin/kin-openapi`
  - [ ] Agregar anotaciones Swagger a handlers
  - [ ] Generar documentación OpenAPI
  - [ ] Configurar endpoint `/api/docs` para visualizar documentación
- [ ] Crear README con instrucciones de setup

### Fase 7: Optimizaciones

**⚠️ ANTES DE CONTINUAR**: Confirmar que la Fase 6 está completa y la aplicación está probada y documentada.
- [ ] Agregar paginación a endpoints de listado
- [ ] Agregar filtros y búsqueda
- [ ] Optimizar queries con índices
- [ ] Implementar caché si es necesario
- [ ] Optimizar inyección de dependencias (usar interfaces para mejor testabilidad)

---

## Consideraciones de Implementación

### Multi-Tenancy
- El `tenant_id` debe extraerse del contexto en cada request (mediante middleware)
- Todas las queries deben filtrar por `tenant_id` para asegurar aislamiento de datos
- Validar que las operaciones entre sucursales pertenezcan al mismo tenant

### Stock Management
- El stock es **único por tenant**, no por sucursal
- Al listar productos, mostrar el stock total del tenant
- Los traspasos no afectan el stock total, solo la distribución entre sucursales
- Implementar reservas de stock durante traspasos pendientes

### Validaciones Clave
- No permitir eliminar categorías que tengan productos asociados
- No permitir traspasos entre sucursales de diferentes tenants
- Validar que las cantidades de stock no sean negativas
- Validar que los traspasos tengan stock suficiente disponible

### Dependencias Sugeridas (Idiomáticas en Go)

```go
// Router HTTP (recomendado: chi por su simplicidad y estándares)
github.com/go-chi/chi/v5

// Base de datos PostgreSQL (pgx es más idiomático que GORM)
github.com/jackc/pgx/v5
github.com/jackc/pgx/v5/pgxpool

// Validación
github.com/go-playground/validator/v10

// Variables de entorno
github.com/joho/godotenv

// Logger estructurado (zap es más performante)
go.uber.org/zap

// UUID
github.com/google/uuid

// Migraciones (opcional, pero recomendado)
github.com/golang-migrate/migrate/v4

// Swagger/OpenAPI
github.com/swaggo/swag
github.com/swaggo/http-swagger
// O alternativamente:
github.com/getkin/kin-openapi
github.com/getkin/kin-openapi/openapi3

// JWT Authentication
github.com/golang-jwt/jwt/v5
```

### Convenciones de Código Go

- **Nombres de paquetes**: En minúsculas, una sola palabra, descriptivos
- **Interfaces**: Nombres que terminen en `-er` cuando sea apropiado (ej: `Repository`, `Handler`)
- **Errores**:
  - Errores del dominio en cada paquete `domain/{entity}/entities/errors.go`
  - Usar `errors.New()` o `fmt.Errorf()` con contexto
  - Variables de error: `var ErrCategoryNotFound = errors.New("category not found")`
- **Repositorios**:
  - Archivos sin apellido `_repository`: `repository/category.go` implementa `CategoryRepository`
  - Nombres descriptivos que indican la entidad, no la tecnología
- **Handlers REST**:
  - Un paquete por recurso: `rest/category/`, `rest/product/`, etc.
  - Un archivo por endpoint **SIN apellidos** (el paquete ya indica el recurso):
    - `list.go`: Listar recursos (GET collection) - Ej: `rest/category/list.go`, `rest/product/list.go`
    - `get_by_id.go`: Obtener recurso por ID (GET resource) - Ej: `rest/category/get_by_id.go`, `rest/product/get_by_id.go`
    - `create.go`: Crear recurso (POST) - Ej: `rest/category/create.go`, `rest/product/create.go`
    - `update.go`: Actualizar recurso completo (PUT) - Ej: `rest/category/update.go`, `rest/product/update.go`
    - `parcial_update.go`: Actualizar campos específicos (PATCH) - Ej: `rest/category/parcial_update.go`, `rest/product/parcial_update.go`
    - `remove.go`: Eliminar recurso (DELETE) - Ej: `rest/category/remove.go`, `rest/product/remove.go`
    - `{action}.go`: Acciones específicas - Ej: `rest/transfer/complete.go`, `rest/transfer/cancel.go`, `rest/stock/adjust.go`
  - Entidades de request/response en `rest/{resource}/entities/{resource}.go` dentro de cada paquete
  - Cada archivo contiene un solo handler function con nombre descriptivo
  - **NO usar apellidos** en nombres de archivos (el paquete ya identifica el recurso)
  - El archivo `parcial_update.go` maneja actualizaciones parciales (PATCH) de campos específicos
- **Context**: Pasar `context.Context` como primer parámetro en funciones que hagan I/O
- **Validaciones**: En la capa de dominio, no en los handlers
- **Métodos Privados y Públicos**:
  - Métodos públicos (mayúscula inicial): Solo lo que necesita ser exportado del paquete
  - Métodos privados (minúscula inicial): Funciones auxiliares, helpers, validaciones internas
  - Regla: Si no se usa fuera del paquete, debe ser privado
  - Ejemplo: `func (s *Service) Create()` es público, `func (s *Service) validate()` es privado
- **Diseño de Métodos**:
  - **Métodos pequeños con única responsabilidad**: Cada método debe hacer una sola cosa bien
  - Ideal: Métodos de 10-20 líneas máximo, si excede, considerar dividirlo
  - Un método = una responsabilidad = un nivel de abstracción
  - Si un método tiene múltiples "y" en su descripción, probablemente hace demasiado
  - Ejemplo: `validateCategory()` y `saveCategory()` en lugar de `validateAndSaveCategory()`
- **Parámetros de Métodos**:
  - **Máximo 3 parámetros** por método (según estándares de Go)
  - Si se necesitan más parámetros, usar structs de configuración
  - Ejemplo: `func Create(ctx context.Context, req CreateRequest) error` en lugar de `func Create(ctx, name, description, tenantID)`
  - Structs de configuración hacen el código más legible y mantenible
- **Código Testeable**:
  - **Diseñar para ser testeable**: Métodos pequeños y con responsabilidad única son más fáciles de testear
  - Inyectar dependencias: No crear dependencias dentro de métodos, recibirlas por constructor
  - Evitar dependencias globales: Dificultan el testing
  - Usar interfaces: Permiten crear mocks fácilmente
  - Separar lógica de negocio de I/O: Facilita testing unitario
  - Ejemplo: `func (s *Service) CalculateTotal(stock []Stock) int` es fácil de testear sin DB
- **Comentarios de Código**:
  - **NO agregar comentarios de contexto** que repitan lo que el código ya dice
  - Los comentarios deben explicar el "por qué", no el "qué"
  - Documentar funciones públicas con comentarios descriptivos (go doc)
  - Evitar comentarios obvios como `// Incrementa el contador` sobre `counter++`
  - Comentar solo lógica compleja o decisiones de negocio no obvias
- **Configuración Parametrizable**:
  - **SIEMPRE** usar `config/config.json` para valores parametrizables
  - NO hardcodear valores que puedan cambiar (timeouts, límites, URLs, etc.)
  - Valores que deben ir en config:
    - Timeouts y timeouts de conexión
    - Límites de paginación (default, máximo)
    - Tamaños de pool de conexiones
    - URLs de servicios externos
    - Límites de validación (longitud máxima de strings, etc.)
    - Configuración de logging (niveles, formatos)
  - Secrets y credenciales: Variables de entorno (`.env`), nunca en config.json
  - Ejemplo: `max_page_size: 100` en config, no `const MaxPageSize = 100` en código
- **Idiomaticidad**:
  - **SIEMPRE** seguir las convenciones idiomáticas de Go
  - Preferir código claro y explícito sobre "clever" o complejo
  - Documentar funciones públicas con comentarios descriptivos
  - Manejar errores explícitamente, nunca ignorarlos con `_`
- **Pre-commit y Calidad de Código**:
  - Implementar pre-commit hooks para validar código antes de commit
  - Ejecutar automáticamente: `go fmt`, `golangci-lint`, `go vet`, `go test`
  - No permitir commits con código que no pase las validaciones
  - Ver sección "Pre-commit Setup" para implementación detallada
- **Gestión de Commits**:
  - **Hacer commits frecuentes**: Cada funcionalidad completa o cambio significativo
  - Commits pequeños y atómicos: Un commit = una funcionalidad/corrección
  - Mensajes descriptivos: "feat: add category creation handler" en lugar de "update"
  - No acumular cambios grandes en un solo commit
  - Commits frecuentes facilitan rollback y debugging
- **Inyección de Dependencias**:
  - Servicios reciben repositorios y config por constructor: `NewCategoryService(repo, config)`
  - Handlers reciben servicios y config por constructor: `NewCategoryHandler(service, config)`
- **Configuración**:
  - Cargar desde JSON en `config/config.json`
  - Secrets desde variables de entorno (`.env`)
- **Tests**: Archivos `*_test.go` en el mismo paquete

---

## Ejemplos de Requests/Responses REST

### Crear Categoría

**Request:**
```http
POST /api/v1/categories HTTP/1.1
Host: api.example.com
X-Tenant-ID: 123e4567-e89b-12d3-a456-426614174000
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "name": "Motocicletas",
  "description": "Categoría para motocicletas"
}
```

**Response (201 Created):**
```json
{
  "id": "789e0123-e89b-12d3-a456-426614174002",
  "name": "Motocicletas",
  "description": "Categoría para motocicletas",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Listar Productos con Filtros

**Request:**
```http
GET /api/v1/products?store_id=456e7890-e89b-12d3-a456-426614174001&category_id=789e0123-e89b-12d3-a456-426614174002&page=1&limit=20 HTTP/1.1
Host: api.example.com
X-Tenant-ID: 123e4567-e89b-12d3-a456-426614174000
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "abc123",
      "store_id": "456e7890-e89b-12d3-a456-426614174001",
      "category_id": "789e0123-e89b-12d3-a456-426614174002",
      "name": "Honda CBR 600",
      "description": "Motocicleta deportiva",
      "sku": "HONDA-CBR-600",
      "price": 15000.00,
      "stock": {
        "quantity": 13,
        "reserved_quantity": 2
      },
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

### Actualizar Stock

**Request:**
```http
PUT /api/v1/products/abc123/stock HTTP/1.1
Host: api.example.com
X-Tenant-ID: 123e4567-e89b-12d3-a456-426614174000
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "quantity": 15
}
```

**Response (200 OK):**
```json
{
  "id": "abc123",
  "quantity": 15,
  "reserved_quantity": 2,
  "available_quantity": 13,
  "updated_at": "2024-01-15T11:00:00Z"
}
```

### Crear Traspaso

**Request:**
```http
POST /api/v1/transfers HTTP/1.1
Host: api.example.com
X-Tenant-ID: 123e4567-e89b-12d3-a456-426614174000
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "product_id": "abc123",
  "from_store_id": "456e7890-e89b-12d3-a456-426614174001",
  "to_store_id": "789e0123-e89b-12d3-a456-426614174003",
  "quantity": 3,
  "notes": "Traspaso por inventario"
}
```

**Response (201 Created):**
```json
{
  "id": "transfer-123",
  "product_id": "abc123",
  "from_store_id": "456e7890-e89b-12d3-a456-426614174001",
  "to_store_id": "789e0123-e89b-12d3-a456-426614174003",
  "quantity": 3,
  "status": "pending",
  "notes": "Traspaso por inventario",
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

### Completar Traspaso

**Request:**
```http
PATCH /api/v1/transfers/transfer-123/complete HTTP/1.1
Host: api.example.com
X-Tenant-ID: 123e4567-e89b-12d3-a456-426614174000
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "id": "transfer-123",
  "status": "completed",
  "updated_at": "2024-01-15T11:05:00Z"
}
```

### Error Response (400 Bad Request)

**Response:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Error de validación",
    "details": [
      {
        "field": "quantity",
        "message": "La cantidad debe ser mayor a 0"
      }
    ]
  }
}
```

---

## Notas Finales y Mejores Prácticas

### Arquitectura
- **Scaffolding inicial**: Enfocado en categorías y productos, extensible para futuras funcionalidades
- **Arquitectura Hexagonal**: Permite fácil extensión y testing sin acoplamiento
- **Orientación al Dominio**: Cada agregado tiene su propio paquete con responsabilidades claras

### Multi-Tenancy
- El `tenant_id` se extrae del header `X-Tenant-ID` mediante middleware
- Todas las queries filtran automáticamente por `tenant_id` para aislamiento de datos
- Validar que operaciones entre sucursales pertenezcan al mismo tenant

### Autenticación y Seguridad
- **JWT (JSON Web Tokens)**: Implementar autenticación basada en tokens
- **Middleware de Autenticación**: Validar token JWT en cada request protegido
- **Secret Key**: Almacenar clave secreta para firmar tokens en variables de entorno
- **Expiración de Tokens**: Configurar tiempo de expiración (recomendado: 1 hora)
- **Protección de Endpoints**: Todos los endpoints excepto `/auth/login` requieren token válido
- **Validación de Tenant**: Combinar validación JWT con validación de tenant para doble seguridad

### Stock Management
- **Stock único por tenant**: No por sucursal
- Al listar productos, mostrar stock total del tenant
- Los traspasos no afectan el stock total, solo la distribución entre sucursales
- Implementar reservas de stock durante traspasos en estado 'pending'
- El stock se libera cuando el traspaso se cancela o completa

### Validaciones
- Validaciones de negocio en la capa de dominio
- Validaciones de formato en la capa REST (entidades de request/response)
- No permitir eliminar categorías con productos asociados
- Validar que traspasos tengan stock suficiente disponible
- Validar que cantidades no sean negativas

### Testing
- Tests unitarios para servicios de dominio
- Tests de integración para repositorios
- Tests de integración para handlers HTTP
- Usar mocks para interfaces de repositorio

### Escalabilidad
- La estructura permite agregar nuevos agregados fácilmente
- Cada dominio es independiente y puede evolucionar por separado
- Los repositorios no están separados por tecnología (fácil cambio de DB)
- La capa REST puede evolucionar sin afectar el dominio
- Inyección de dependencias permite mejor testabilidad y mantenimiento
- Configuración centralizada en JSON facilita cambios sin recompilar

### Inyección de Dependencias y Configuración

**Ejemplo de estructura de configuración (`config/config.json`):**
```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0",
    "read_timeout": "30s",
    "write_timeout": "30s"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "max_connections": 100,
    "max_idle_connections": 10,
    "conn_max_lifetime": "5m"
  },
  "pagination": {
    "default_limit": 20,
    "max_limit": 100
  },
  "validation": {
    "max_name_length": 255,
    "max_description_length": 1000
  },
  "logging": {
    "level": "info",
    "format": "json"
  }
}
```

**Nota**: Todos los valores parametrizables deben estar aquí. NO hardcodear en el código.

**Ejemplo de constructor de servicio:**
```go
// domain/category/service.go
type Service struct {
    repo   Repository
    config *Config
    logger logger.Logger
}

func NewService(repo Repository, config *Config, logger logger.Logger) *Service {
    return &Service{
        repo:   repo,
        config: config,
        logger: logger,
    }
}
```

---

## 🤖 Guía para Vibe Coding (Desarrollo con IA)

### Principios para Maximizar la Efectividad del Desarrollo con IA

#### 1. Comunicación Clara y Específica
- **Ser explícito**: Especificar exactamente qué se necesita, no asumir que la IA "entiende el contexto"
- **Proporcionar contexto**: Cuando se solicita un cambio, incluir información relevante (archivo, línea, error, etc.)
- **Iterar incrementalmente**: Pedir cambios pequeños y específicos, no grandes refactorizaciones de una vez
- **Validar resultados**: Revisar el código generado antes de aceptarlo

#### 2. Estructura y Organización
- **Seguir el documento**: La IA debe seguir estrictamente `PROMPT_CURSOR.md`
- **Una tarea a la vez**: Solicitar implementación de una funcionalidad completa antes de pasar a la siguiente
- **Confirmar avances**: Pedir confirmación explícita antes de continuar con la siguiente fase
- **Mantener consistencia**: Asegurar que cada nueva implementación siga los mismos patrones establecidos

#### 3. Mejores Prácticas para Solicitudes
- **Formato de solicitudes**:
  ```
  ✅ BUENO: "Implementa el handler create.go para categorías siguiendo la estructura del documento"
  ❌ MALO: "Haz el create de categorías"
  ```
- **Incluir referencias**: Cuando se pide algo similar a lo existente, mencionar el archivo de referencia
- **Especificar validaciones**: Si hay reglas de negocio específicas, mencionarlas explícitamente
- **Pedir explicaciones**: Si algo no está claro, pedir que la IA explique su enfoque antes de implementar

#### 4. Revisión y Validación
- **Revisar código generado**: Siempre revisar el código antes de aceptarlo
- **Validar idiomaticidad**: Verificar que el código siga convenciones de Go
- **Probar funcionalidad**: Ejecutar tests y verificar que funciona como se espera
- **Solicitar mejoras**: Si algo no cumple con las expectativas, pedir correcciones específicas

#### 5. Gestión de Errores y Problemas
- **Reportar errores completos**: Incluir mensaje de error, stack trace, y contexto
- **Solicitar explicaciones**: Si algo no funciona, pedir que la IA explique qué puede estar mal
- **Iterar sobre soluciones**: Si una solución no funciona, proporcionar feedback y pedir alternativa
- **Documentar decisiones**: Si se toma una decisión diferente a la documentada, actualizar el documento

#### 6. Optimización del Flujo de Trabajo
- **Usar el documento como referencia**: La IA debe leer `PROMPT_CURSOR.md` cuando haya dudas
- **Mantener el contexto**: Si se pierde contexto, pedir que la IA relea el documento completo
- **Agrupar tareas relacionadas**: Solicitar implementación de funcionalidades relacionadas juntas
- **Revisar antes de continuar**: No avanzar a la siguiente fase sin validar la actual
- **Commits frecuentes**: Hacer commit después de cada funcionalidad completa
  - No acumular muchos cambios sin commit
  - Commits pequeños y atómicos facilitan rollback
  - Mensajes descriptivos: "feat: add category creation handler"

#### 7. Comandos Útiles para la IA
- "Lee el archivo PROMPT_CURSOR.md completo antes de continuar"
- "Implementa [funcionalidad] siguiendo exactamente la estructura del documento"
- "Valida que el código sea idiomático en Go"
- "Confirma que todos los valores parametrizables están en config.json"
- "Verifica que los métodos privados/públicos estén correctamente definidos"
- "Revisa que no hay comentarios de contexto innecesarios"

#### 8. Checklist Pre-Implementación
Antes de pedir implementación, verificar:
- [ ] ¿Está claro qué se necesita implementar?
- [ ] ¿Se ha especificado el archivo/paquete donde debe ir?
- [ ] ¿Se han mencionado las validaciones o reglas de negocio?
- [ ] ¿Se ha referenciado el documento PROMPT_CURSOR.md?
- [ ] ¿Se ha confirmado que la fase anterior está completa?

#### 9. Checklist Post-Implementación
Después de recibir código, verificar:
- [ ] ¿El código sigue las convenciones idiomáticas de Go?
- [ ] ¿Los métodos públicos/privados están correctamente definidos?
- [ ] ¿Los métodos son pequeños con única responsabilidad?
- [ ] ¿Los métodos tienen máximo 3 parámetros (o usan structs)?
- [ ] ¿El código es fácil de testear (inyección de dependencias, interfaces)?
- [ ] ¿Los valores parametrizables están en config.json?
- [ ] ¿No hay comentarios de contexto innecesarios?
- [ ] ¿El código compila sin errores?
- [ ] ¿Pasa los hooks de pre-commit (go fmt, go vet, golangci-lint)?
- [ ] ¿Sigue la estructura del documento?
- [ ] ¿Hacer commit después de validar todo lo anterior?

---

## 🔄 Instrucciones de Continuidad

### Si el Contexto se Agota o se Pierde Información

**IMPORTANTE**: Si durante la implementación:
- El contexto de la conversación se agota
- Se pierde información sobre decisiones tomadas
- Hay dudas sobre la estructura o convenciones
- Se necesita recordar especificaciones del proyecto

**ACCIÓN REQUERIDA**:
1. **LEER COMPLETAMENTE** el archivo `PROMPT_CURSOR.md` desde el inicio
2. Verificar que todas las decisiones sigan las especificaciones del documento
3. Asegurar que la estructura de archivos coincida con la documentada
4. Validar que las convenciones de naming sean consistentes
5. Confirmar que la arquitectura hexagonal se mantiene
6. Revisar que la conexión con Supabase sigue las instrucciones detalladas

### Recordatorios Clave

- ✅ **SIEMPRE** pedir confirmación antes de pasar a la siguiente fase
- ✅ **SIEMPRE** seguir convenciones idiomáticas de Go
- ✅ **SIEMPRE** validar conexión con Supabase antes de continuar
- ✅ **SIEMPRE** leer el `.md` completo si hay dudas o se pierde contexto
- ✅ **NUNCA** hardcodear credenciales o información sensible
- ✅ **NUNCA** avanzar sin confirmación explícita del usuario

---

**Fin del Documento**
