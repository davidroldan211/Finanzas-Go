# SPEC 01 — Registro con correo

> **Estado:** Implementado
> **Tipo:** feat
> **Depende de:** —
> **Historias:** HU-01
> **Capacidades:** CAP-01
> **Fecha:** 2026-11-14
> **Objetivo:** Que un visitante cree su usuario con correo, nombre, apellido y contraseña mediante un registro que no revela qué correos existen y exige contraseñas robustas.

## 1. Por qué existe esta spec

El registro ya existe: `POST /api/v1/users` en `internal/users`. Funciona, pero no cumple lo que F3 pide a un "registro público seguro":

- Un correo ya registrado responde `409 El email ya está registrado`. Cualquiera puede averiguar si un correo tiene usuario (A5).
- La contraseña mínima es de 6 caracteres y no tiene máximo. bcrypt ignora todo lo que pase de 72 bytes sin avisar.
- El dominio acepta cualquier correo con un `@`.
- El correo de un usuario eliminado (soft delete) pasa el chequeo de `ExistsByEmail`, pero choca con el `UNIQUE` de la columna y responde 500.
- No hay forma de saber si el usuario confirmó su correo, lo que HU-02 necesita.

Esta spec endurece el registro existente y deja la base para HU-02. No crea un registro nuevo ni cambia la ruta.

## 2. Alcance

**Entra:**

- `POST /api/v1/users` responde siempre `202 Accepted` con el mismo mensaje genérico cuando los datos son válidos, exista o no el correo.
- Si el correo ya existe (activo o eliminado), no se crea nada y la respuesta es idéntica a la de un registro nuevo.
- La contraseña se hashea en ambos caminos, antes de consultar si el correo existe, para igualar el tiempo de respuesta.
- Una carrera entre dos registros con el mismo correo termina en `202` para ambos, no en 500.
- Política de contraseñas en el dominio: 12 caracteres mínimo, 72 bytes máximo, no estar en la lista de contraseñas comunes, no contener la parte local del correo, el nombre ni el apellido.
- Validación del correo con `net/mail` en el dominio.
- Columna nueva `users.email_verified_at`, que el registro deja en `NULL`.
- El registro sigue pidiendo `email`, `first_name`, `last_name` y `password`, y sigue creando siempre rol `user` (H1).

**Fuera de alcance (para specs futuras):**

- Generar, enviar y confirmar el código de verificación (HU-02).
- Bloquear el inicio de sesión mientras el correo no esté verificado (HU-02).
- Avisar por correo al dueño de un correo ya registrado cuando alguien intenta registrarse con él (HU-02, cuando exista el envío de correo).
- Límite de peticiones al registro y al resto de rutas de identidad (spec transversal propia, A4 y A16).
- Autorización de tratamiento de datos de la Ley 1581 (antes de F6, junto con la política de tratamiento).
- Moneda principal y zona horaria por defecto (HU-20, CAP-04).
- Reusar el correo de un usuario eliminado (llega con el borrado real de CAP-05).
- Aplicar la política de contraseñas a la recuperación y al cambio de contraseña (HU-04).

## 3. Modelo de datos

**Migración** — `db/migrations/00004_add_email_verified_at_to_users.sql` (nueva; las migraciones 00001–00003 no se tocan).

```sql
-- +goose Up
ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMPTZ NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
```

Los usuarios que ya existan quedan con `NULL` (no verificados). No hay usuarios reales todavía.

**Dominio** — `internal/users/domain/user.go`

```go
type User struct {
	// ...campos actuales...
	EmailVerifiedAt *time.Time // nil = correo sin confirmar
}
```

`NewUser` deja `EmailVerifiedAt` en `nil`. Ningún caso de uso de esta spec lo cambia.

**Dominio** — `internal/users/domain/password_policy.go`

```go
const (
	PasswordMinChars = 12 // utf8.RuneCountInString
	PasswordMaxBytes = 72 // límite de bcrypt
)

// ValidatePassword devuelve el mensaje del campo "password", o "" si es válida.
func ValidatePassword(plain, email, firstName, lastName string) string
```

Reglas, en este orden (se informa solo la primera que falle):

1. Menos de 12 caracteres: `"la contraseña debe tener al menos 12 caracteres"`.
2. Más de 72 bytes: `"la contraseña es demasiado larga"`.
3. En minúsculas, está en la lista de contraseñas comunes: `"esa contraseña es muy común"`.
4. En minúsculas, contiene la parte local del correo, el nombre o el apellido, cuando ese valor tiene 3 caracteres o más: `"la contraseña no puede contener tu correo ni tu nombre"`.

No hay reglas de composición (mayúsculas, símbolos).

**Lista de contraseñas comunes** — `internal/users/domain/common_passwords.txt`, embebida con `//go:embed`.

- Top 10.000 de SecLists (`Passwords/Common-Credentials/10k-most-common.txt`, licencia MIT), en minúsculas.
- Más una lista corta escrita a mano de contraseñas comunes en español y Colombia.
- Una por línea. Se carga una vez en un `map[string]struct{}` al iniciar el paquete.
- La cabecera del archivo cita la fuente y la licencia.

**Validación de correo** — `validateUserFields` en `internal/users/domain/user.go` reemplaza el chequeo de `@`:

- `net/mail.ParseAddress(email)` sin error.
- La dirección parseada es igual a la entrada normalizada (rechaza `"Ana <a@b.co>"`).
- La parte de dominio contiene un punto.
- 254 caracteres como máximo.

**Puerto de entrada** — `internal/users/port/in/user_service.go`

```go
// Create es el registro público. Devuelve nil tanto si creó el usuario
// como si el correo ya existía: el llamador no puede distinguirlos.
Create(ctx context.Context, cmd CreateUserCommand) error
```

**Adaptador de salida** — `internal/users/adapter/out/postgres/`

- `userModel` agrega `EmailVerifiedAt *time.Time`, con su mapeo en `toDomain` y `toModel`.
- `ExistsByEmail` cuenta también las filas eliminadas (`Unscoped()`, sin el filtro `deleted_at IS NULL`).
- `Save` traduce una violación de unicidad de Postgres (`*pgconn.PgError` con código `23505`) a `domain.ErrEmailTaken`.

**Adaptador de entrada** — `internal/users/adapter/in/http/user_dto.go`

- `CreateUserRequest` deja solo `binding:"required"` en los cuatro campos. Se quitan `email` y `min=6`: el dominio es la única fuente de las reglas.
- Respuesta `202`: `{"message": "Registro recibido. Revisa tu correo para confirmarlo."}`. Sin `user`, sin `id`.

## 4. Impacto en seguridad

- **Datos que toca:** correo, nombre y apellido (Confidencial) y la contraseña (Restringido). Ver SEC-0 §1. `email_verified_at` es Confidencial.
- **Quién puede acceder:** `POST /api/v1/users` es público (visitante). La respuesta no devuelve ningún dato del usuario. Las demás rutas de `users` no cambian.
- **Escenarios de SEC-0:**
  - A1 (crearse administrador): el registro sigue ignorando `role`; esta spec no lo reabre.
  - A5 (averiguar si un correo existe): respuesta, código y tiempo iguales en ambos caminos.
  - A4 (contraseñas filtradas): la lista común y el mínimo de 12 reducen las contraseñas adivinables.
  - A14 (filtrar información interna): la violación de unicidad deja de producir un 500.
- **Controles:**
  - C1. `application.Create` devuelve `nil` cuando el correo ya existe, sin llamar a `Save`. El handler responde `202` con el mismo cuerpo.
  - C2. `application.Create` llama a `PasswordHasher.Hash` antes de `ExistsByEmail`, en todos los caminos que pasan la validación.
  - C3. `ExistsByEmail` considera existentes los correos de usuarios eliminados.
  - C4. `Save` traduce el código `23505` a `ErrEmailTaken`, y `application.Create` lo convierte en `nil`.
  - C5. `domain.ValidatePassword` aplica las cuatro reglas de la sección 3, y `application.Create` la llama antes de hashear.
  - C6. La respuesta de registro no incluye contraseña, hash, `id` ni ningún otro campo del usuario.
  - C7. Ni `application` ni el handler registran en logs el correo del intento cuando el correo ya existe.
  - C8. El registro crea siempre rol `user`, aunque el cuerpo traiga `role` (regresión de H1).
- **Datos nuevos que se piden:** ninguno. Los cuatro campos ya se pedían y están permitidos por SEC-0 §2.
- **Dependencias nuevas:** ninguna de Go. `pgconn` ya viene con `pgx/v5`. La lista de SecLists es un archivo de datos con licencia MIT, citado en su cabecera.

## 5. Plan de implementación

1. **Migración** — Crear `db/migrations/00004_add_email_verified_at_to_users.sql` con el `Up` y el `Down` de la sección 3. Verificación: `make migrate-up`, `make migrate-down` y `make migrate-up` de nuevo funcionan contra el Postgres local.

2. **Dominio: correo y verificación** — Agregar `EmailVerifiedAt` a `User` y reemplazar la validación de correo en `validateUserFields` por la de `net/mail`. Pruebas en `user_test.go`: `"Ana <a@b.co>"`, `"a@b"`, `"sin-arroba"` y un correo de 255 caracteres dan `fields.email`; `"Ana@Ejemplo.co "` se normaliza a `"ana@ejemplo.co"`; `NewUser` deja `EmailVerifiedAt` en `nil`.

3. **Dominio: política de contraseñas** — Crear `password_policy.go` y `common_passwords.txt` con `//go:embed`. Pruebas en `password_policy_test.go`: 11 caracteres falla; 12 pasa; 73 bytes falla; 18 caracteres multibyte (`"ñ"`, 36 bytes) pasa; una contraseña de la lista falla aunque venga en mayúsculas; una que contiene la parte local del correo, el nombre o el apellido falla; un nombre de 2 caracteres contenido en la contraseña no la hace fallar.

4. **Aplicación: orden del registro** — En `user_service.go`, `Create` valida con `domain.NewUser` y `domain.ValidatePassword` (fusionando los campos en un solo `ValidationError`), luego hashea, luego consulta `ExistsByEmail`. La firma no cambia todavía. Pruebas con stubs en `user_service_test.go`: datos inválidos devuelven `ValidationError` sin llamar al hasher; correo existente llama al hasher y no llama a `Save` (C2); rol del usuario guardado es `user` (C8).

5. **Puerto y aplicación: registro indistinguible** — Cambiar `in.UserService.Create` a `Create(ctx, cmd) error`. `Create` devuelve `nil` si el correo existe y si `Save` devuelve `ErrEmailTaken`. Actualizar el handler a `202` con el cuerpo de la sección 3 y el stub de `user_handler_test.go`. Pruebas de aplicación: correo existente devuelve `nil` sin `Save` (C1); `Save` con `ErrEmailTaken` devuelve `nil` (C4); ningún camino escribe en el logger (C7). Pruebas de handler: registro nuevo y correo existente dan `202` con cuerpos idénticos byte a byte; el cuerpo no contiene `password`, `id` ni el correo (C6); `TestCreateUser_EmailTaken_Returns409` se reemplaza por la prueba equivalente con `202`.

6. **DTO: una sola fuente de reglas** — Quitar `email` y `min=6` de `CreateUserRequest`. Prueba de handler: una contraseña de 6 caracteres con un servicio real (stubs de repositorio y hasher) da `422` con `fields.password` del dominio (C5).

7. **Adaptador de salida** — En `adapter/out/postgres`: agregar `EmailVerifiedAt` a `userModel` y a los mappers; `ExistsByEmail` con `Unscoped()` (C3); `Save` traduce `23505` con un helper `isUniqueViolation(err) bool` (C4). Pruebas: ida y vuelta de `EmailVerifiedAt` en `user_model_test.go`; `isUniqueViolation` reconoce `&pgconn.PgError{Code: "23505"}` envuelto con `fmt.Errorf("%w")` y rechaza otro código.

## 6. Criterios de aceptación

- [x] `make ci` pasa.
- [x] `make arch` pasa: `domain` solo usa `net/mail`, `embed` y la librería estándar.
- [x] `make migrate-up` y `make migrate-reset` funcionan con la migración 00004. (job Migraciones)
- [x] Registro con un correo nuevo responde `202` y crea el usuario con `email_verified_at` en `NULL`. (prueba de aplicación y prueba manual contra Postgres local)
- [x] Registro con un correo ya registrado responde `202` con el mismo cuerpo byte a byte que uno nuevo y no crea ninguna fila. (C1; prueba de handler y de aplicación)
- [x] Con un correo existente, el hasher se llama una vez antes de `ExistsByEmail`. (C2; prueba de aplicación)
- [x] Registro con el correo de un usuario eliminado responde `202` y no da 500. (C3; prueba manual contra Postgres local)
- [x] `Save` con violación `23505` termina en `202`. (C4; prueba de `isUniqueViolation` y de aplicación)
- [x] Una contraseña de 11 caracteres, una de 73 bytes, una de la lista común y una que contiene el nombre responden `422` con `fields.password`. (C5; pruebas de dominio y de handler)
- [x] La respuesta `202` no contiene `password`, `id`, el correo ni el hash. (C6; prueba de handler)
- [x] Ningún camino de `Create` escribe el correo en logs. (C7; prueba de aplicación y revisión del diff)
- [x] Un cuerpo con `"role":"admin"` crea un usuario con rol `user`. (C8; prueba de aplicación)
- [x] `common_passwords.txt` cita en su cabecera la fuente SecLists y la licencia MIT.

## 7. Decisiones

- **Sí:** endurecer `POST /api/v1/users` en vez de crear una ruta nueva. El caso de uso crea un usuario y vive en `users`; no hay frontend que dependa de la ruta.
- **No:** mover el registro a `/api/v1/auth/register`. `auth` tendría que escribir en la tabla de `users` o depender de su servicio, y rompe el patrón actual.
- **Sí:** `202 Accepted` uniforme con mensaje genérico. Evita la enumeración de correos (A5).
- **No:** mantener `409`. Mejor experiencia, pero revela qué correos están registrados.
- **Sí:** `Create` devuelve solo `error`. El handler no puede filtrar si el correo existía porque no recibe esa información.
- **Sí:** hashear antes de consultar el correo. Iguala el costo de bcrypt en ambos caminos sin dependencias nuevas.
- **Sí:** columna `email_verified_at` ahora, sin bloquear el login. HU-02 decide el bloqueo cuando exista el envío del código; bloquear antes dejaría a todos sin poder entrar.
- **No:** bloquear el login en esta spec, ni dejar todo el estado de verificación para HU-02.
- **Sí:** mantener `first_name` y `last_name` obligatorios. Sin migración de columnas; son datos de identificación permitidos.
- **No:** pedir solo correo y contraseña. Exigiría volver opcionales las columnas y completar el nombre en el perfil.
- **Sí:** 12 caracteres mínimo (OWASP ASVS 2.1.1), 72 bytes máximo (límite de bcrypt), sin reglas de composición.
- **No:** mínimo de 8 (NIST 800-63B). Menos fricción, menos fuerza.
- **Sí:** lista común de SecLists 10k más contraseñas en español escritas a mano, embebida en el binario. Sin servicio externo ni dependencia de Go.
- **Sí:** rechazar contraseñas que contienen el correo o el nombre (ASVS 2.1.7), con umbral de 3 caracteres para no rechazar por nombres muy cortos.
- **Sí:** validar el correo con `net/mail` en el dominio. Biblioteca estándar; HU-02 prueba igual que el correo es del usuario.
- **Sí:** tratar el correo de un usuario eliminado como existente. Corrige el 500 actual sin migración. Reusar correos llega con el borrado real (CAP-05).
- **No:** índice único parcial `WHERE deleted_at IS NULL`. Cambia la semántica del login por correo y adelanta una decisión de CAP-05.
- **Sí:** el cambio en `ExistsByEmail` aplica también a `Update`. Un usuario no puede cambiar su correo al de un usuario eliminado, lo que hoy también terminaría en 500.
- **Sí:** la política de contraseñas vive en `internal/users/domain`. Si HU-04 la necesita desde otro módulo, esa spec decide si se mueve a `internal/shared`.
- **No:** limitador de peticiones en esta spec. Sirve también a login, verificación y recuperación; va en su propia spec transversal.
- **No:** autorización de tratamiento de datos en esta spec. No existe política de tratamiento escrita; la normativa es requisito de entrada a SEC-6.

## 8. Riesgos

| Riesgo | Mitigación |
| --- | --- |
| Una persona que olvidó que ya tenía usuario se registra de nuevo y no entiende por qué no puede entrar. | El mensaje genérico la manda a su correo. Con HU-02 recibe un aviso de que ya tiene usuario. Mientras tanto, puede usar la recuperación de contraseña (HU-04). |
| La diferencia de tiempo no queda del todo igualada (el `INSERT` solo ocurre en un camino). | El costo dominante es bcrypt (~100 ms frente a ~1 ms del `INSERT`). El limitador de peticiones reducirá el número de mediciones posibles. |
| La lista común agrega ~80 KB al binario. | Aceptable para el beneficio. Se carga una vez al iniciar. |
| Sin limitador, alguien puede crear muchas cuentas. | Riesgo aceptado hasta la spec del limitador (A16). |

## Lo que **no** está en esta spec

- Envío y confirmación del código de verificación (HU-02).
- Bloqueo del login sin correo verificado (HU-02).
- Aviso al dueño de un correo ya registrado (HU-02).
- Limitador de peticiones (spec transversal).
- Autorización de tratamiento de datos de la Ley 1581.
- Moneda principal y zona horaria (HU-20).
- Reuso de correos de usuarios eliminados (CAP-05).
- Política de contraseñas en recuperación y cambio de contraseña (HU-04).

Cada uno, si llega, va en su propia spec.
