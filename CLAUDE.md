# CLAUDE.md

Este archivo da guía a Claude Code (claude.ai/code) al trabajar con código en este repositorio.

## Comandos

```bash
make run       # go run ./cmd/finanzas
make build     # compila binario a bin/
make test      # go test ./...
make lint      # go vet ./... + make arch
make arch      # verifica que domain/application/port no importen gin/gorm/config/httpx
make fmt       # go fmt ./...
make coverage  # test + reporte de cobertura, falla si total < 5%
```

Test individual: `go test ./internal/users/application/... -run TestName -v`

Requiere Docker (`docker compose up -d`) para Postgres, y archivo `.env` (ver `.env.template`) — `config.LoadConfig()` llama `godotenv.Load()` y hace `log.Fatal` si falta `.env`. Las migraciones en `db/migrations/*.sql` se aplican a mano (no hay `AutoMigrate` ni herramienta de migración).

## Arquitectura

Arquitectura hexagonal (puertos y adaptadores) modular en Go, por módulo bajo `internal/<modulo>/`:

```
internal/<modulo>/
  domain/               # entidad pura + sentinelas de error; sin tags json/gorm, sin uuid.Nil como "válido"
  port/in/               # puerto de entrada: interfaz <Modulo>Service + Commands/Queries que consume el adaptador in
  port/out/               # puertos de salida: interfaces que consume application (Repository, PasswordHasher, TokenProvider...)
  application/            # implementa port/in; depende SOLO de port/out — nunca de gin, gorm, config ni httpx
  adapter/in/http/         # handlers Gin + DTOs + routes + traducción de errores de dominio a httpx.AppError
  adapter/out/postgres/    # implementa port/out contra GORM; tiene su propio modelo de DB + mappers toDomain/toModel
  <modulo>_module.go       # composition root del módulo: arma adapter/out -> application -> adapter/in
```

Módulos existentes: `auth`, `users`, `verification`. Regla verificada por `make arch`: **`domain`, `port/*` y `application` nunca importan `gin`, `gorm`, `finanzas-api/config` ni `internal/httpx`** — esas dependencias están confinadas a `adapter/*`. Al agregar una funcionalidad: definir/extender el puerto primero (`port/in` si es una operación nueva expuesta, `port/out` si es una dependencia externa nueva), implementar en `application`, exponer vía `adapter/in/http`.

`cmd/finanzas/main.go` es el composition root general: carga config, abre la DB, construye los adaptadores de seguridad compartidos (`security.NewBcryptHasher`, `security.NewHMACTokenProvider`), construye cada `<modulo>.NewModule(db, ...)` y registra sus rutas. Para agregar un módulo nuevo, seguir este mismo patrón y registrarlo en `main.go`.

**Cómo `auth` lee usuarios sin acoplarse a `users`:** `auth` NO reutiliza el repositorio de `users`. Tiene su propio adaptador de salida read-only (`internal/auth/adapter/out/postgres/user_credentials_repository.go`) con una proyección de 5 columnas sobre la misma tabla `users` — la fuente de verdad del esquema sigue siendo `db/migrations/users.sql`. Es el patrón a replicar si otro módulo necesita leer datos de un módulo distinto: un adaptador de salida propio, nunca importar el paquete de infraestructura ajeno.

**El guard de auth** (`internal/auth/adapter/in/http/auth_guard.go`, tipo `Guard`) recibe un `port/out.TokenProvider` por constructor — no un secreto ni una función suelta. `Module.Guard()` lo expone como `httpx.AuthGuard` (`type AuthGuard func(roles ...string) gin.HandlerFunc`, definido en `internal/httpx/gin.go`), que es lo que otros módulos reciben para proteger sus rutas sin importar `auth` — solo importan `httpx`.

### Piezas transversales

- `config/env.config.go` — config tipada desde variables de entorno (`Database`, `JWT`, `Server`, `App`). Solo `main.go` la importa; ningún módulo la conoce.
- `internal/shared/db/postgres.go` — setup de conexión GORM Postgres (`NewPostgresDB`), usado solo por `main.go`.
- `internal/shared/security/bcrypt_hasher.go` — `BcryptHasher` implementa a la vez `users/port/out.PasswordHasher` (`Hash`) y `auth/port/out.PasswordVerifier` (`Matches`); se construye una sola vez en `main.go` y se inyecta en ambos módulos.
- `internal/shared/security/hmac_token_provider.go` — `HMACTokenProvider` implementa `auth/port/out.TokenProvider` (`Issue`/`Verify`) con HMAC-SHA256 hecho a mano (no es una librería JWT); guarda el secreto y el TTL, valida el header `alg`. Mantiene el tag JSON `user_id` en el payload por compatibilidad con tokens emitidos antes de la migración a hexagonal.
- `internal/httpx/apperr.go` + `internal/httpx/gin.go` — `AppError` (status/code/message/fields) con factories (`httpx.BadRequest`, `httpx.NotFound`, `httpx.Conflict`, `httpx.Validation`, etc.) y su adaptador Gin: `httpx.Abort(c, err)` escribe el envelope `{"error":{"code","message","fields"}}` y nunca filtra el mensaje interno de un error 5xx (solo lo loguea server-side); `httpx.BindJSON(c, &dst)` hace bind + valida, respondiendo 422 con `fields` poblados o 400 si el body no es JSON.

### Ruteo

Todas las rutas se montan bajo `/api/v1` directamente en el `*gin.Engine` compartido en `main.go`. La protección de rutas se aplica por ruta vía `guard("rol1", "rol2")` (tipo `httpx.AuthGuard`) pasado a cada `Setup<Modulo>Routes`/`RegisterRoutes`.

### IDs y errores

- `User.ID` es un `google/uuid.UUID` directo en el dominio (no un tipo propio envuelto) — deliberado: evita el riesgo de que GORM pierda `driver.Valuer` al mapear la PK. Los adaptadores de salida siempre usan `.Where("id = ?", id)` explícito, nunca la condición posicional `First(&m, id)`.
- Los repositorios traducen `gorm.ErrRecordNotFound` a un sentinel de dominio (p.ej. `domain.ErrUserNotFound`) **dentro del adaptador de salida**; `application` y `adapter/in/http` nunca ven un error de GORM directamente.
