# Motico API

API REST en Go para el manejo de inventarios multi-tenant.

## Arquitectura

- **Lenguaje**: Go (Golang)
- **Arquitectura**: Hexagonal (Ports & Adapters)
- **Base de Datos**: Supabase (PostgreSQL)
- **Router**: Chi

## Requisitos

- Go 1.21 o superior
- PostgreSQL (Supabase)
- Pre-commit (opcional, para validaciones)

## Configuración

### 1. Variables de Entorno

Copia `.env.example` a `.env` y configura las variables:

```bash
cp .env.example .env
```

Edita `.env` con tus credenciales de Supabase:
- `DB_HOST`: Host de Supabase
- `DB_PORT`: Puerto (5432)
- `DB_USER`: Usuario
- `DB_PASSWORD`: Contraseña
- `DB_NAME`: Nombre de la base de datos
- `JWT_SECRET_KEY`: Clave secreta para JWT (mínimo 32 caracteres)

### 2. Configuración JSON

Edita `config/config.json` con la configuración de la aplicación.

### 3. Migraciones

Ejecuta las migraciones SQL en tu base de datos Supabase:

```bash
# Ejecuta el contenido de migrations/001_initial_schema.sql en Supabase
```

## Instalación

```bash
# Instalar dependencias
go mod download

# Instalar pre-commit (opcional)
pre-commit install
```

## Ejecución

```bash
go run cmd/api/main.go
```

El servidor iniciará en `http://0.0.0.0:8080`

## Endpoints

### Autenticación

```
POST /api/v1/auth/login
```

## Estructura del Proyecto

```
motico-api/
├── cmd/api/              # Punto de entrada
├── internal/
│   ├── domain/           # Lógica de negocio
│   ├── repository/       # Implementaciones de repositorios
│   └── rest/            # Handlers HTTP
├── pkg/                  # Código reutilizable
├── config/               # Configuración
└── migrations/          # Migraciones SQL
```

## Desarrollo

### Pre-commit Hooks

Los hooks de pre-commit validan automáticamente:
- Formato de código (`go fmt`)
- Análisis estático (`go vet`)
- Tests (`go test`)
- Linter (`golangci-lint`)

### Ejecutar Tests

```bash
go test ./...
```

### Ejecutar Linter

```bash
golangci-lint run ./...
```

## Licencia

MIT
