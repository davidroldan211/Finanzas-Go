# CLAUDE.md

Este archivo da guía a Claude Code (claude.ai/code) al trabajar con código en este repositorio.

## Comandos

```bash
make run       # go run ./cmd/finanzas
make build     # compila binario a bin/
make test      # go test ./...
make lint      # go vet ./...
make fmt       # go fmt ./...
make coverage  # test + reporte de cobertura, falla si total < 5%
```

Test individual: `go test ./internal/auth/handler/... -run TestName -v`

Requiere Docker (`docker compose up -d`) para Postgres, y archivo `.env` (ver `.env.template`) — `config.LoadConfig()` llama `godotenv.Load()` y hace `log.Fatal` si falta `.env`.

## Arquitectura

Clean Architecture modular en Go, por módulo bajo `internal/<modulo>/`. Cada módulo tiene 4 capas más un archivo de wiring:

```
internal/<modulo>/
  domain/      # struct entidad + interfaz Repository + interfaz UseCase (los contratos)
  usecase/     # implementa la interfaz UseCase, depende solo de la interfaz Repository
  repository/  # implementa la interfaz Repository contra GORM/Postgres (tiene su propio modelo de DB, separado de la entidad domain)
  handler/     # handlers HTTP de Gin, depende de la interfaz UseCase
  routes/      # registro de rutas Gin, conectado con el handler
  <modulo>_module.go  # composition root: arma repo -> usecase -> handler para este módulo
```

Módulos existentes: `auth`, `users`, `verification`. `middleware` (guard de auth/rol JWT) vive en `internal/middleware`, no dentro de un módulo.

Dependencias siempre apuntan hacia adentro, hacia las interfaces de `domain` (repository y usecase son interfaces definidas en `domain`, las implementaciones concretas viven en las capas externas) — al agregar una funcionalidad, definir/extender primero la interfaz en `domain`, luego implementar en `usecase`/`repository`, luego exponer vía `handler`/`routes`.

`cmd/finanzas/main.go` es el composition root general de la app: carga config, abre la DB, construye cada `<modulo>.New<Modulo>Module(db, ...)`, luego llama cada `Setup<Modulo>Routes(r, ...)`. Para agregar un módulo nuevo, seguir este mismo patrón de wiring y registrarlo en `main.go`.

Nota: `auth` reutiliza `users/repository` directamente (no tiene repository propio) porque auth solo necesita buscar usuarios por email/password.

### Piezas transversales

- `config/env.config.go` — config tipada cargada desde variables de entorno (secciones `Database`, `JWT`, `Server`, `App`), cada una con defaults vía `getEnv`/`getEnvAsInt`.
- `shared/db/postgres.go` — setup de conexión GORM Postgres (`NewPostgresDB`).
- `shared/security/token.go` — token tipo JWT hecho a mano con HMAC-SHA256 (`GenerateToken`/`ParseToken`), no es una librería JWT.
- `shared/security/password.go` — hashing con bcrypt.
- `internal/middleware/auth_middleware.go` — `Middleware.Handler(roles ...string)` retorna un middleware de Gin que parsea el bearer token y opcionalmente exige pertenencia a rol; setea `userID`/`userRole` en el contexto de Gin.
- `internal/httpx/apperr.go` — tipo `AppError` con status/code/message/fields y funciones factory (`httpx.BadRequest`, `httpx.NotFound`, `httpx.Validation`, etc.) para respuestas de error de API consistentes; usar `httpx.Wrap(err)` como último recurso para convertir a un `AppError` 500.

### Ruteo

Todas las rutas se montan bajo `/api/v1` directamente en el `*gin.Engine` compartido en `main.go` (no hay engine por módulo). La protección de rutas se aplica por ruta vía `authMiddleware("rol1", "rol2")` pasado a `Setup<Modulo>Routes`.

### IDs

`User.ID` es un `google/uuid.UUID`, no un int autoincremental — mantener esto consistente al agregar entidades/repositories que referencien usuarios.
