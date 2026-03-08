---
name: start
description: Arranca el servidor de Motico API en modo desarrollo
disable-model-invocation: true
allowed-tools: Bash
---

Arranca el servidor de Motico API ejecutando:

```bash
go run cmd/api/main.go
```

El servidor inicia en `http://0.0.0.0:8080` con rutas bajo `/api/v1`.

Si el comando falla:
1. Verifica que exista el archivo `.env` con las variables de entorno necesarias
2. Verifica la conexión a la base de datos PostgreSQL
3. Muestra el error al usuario con sugerencias para solucionarlo
