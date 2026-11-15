# Plantilla de spec — Finanzas-Go

Referencia que consulta `/spec` al generar una spec. Cada sección trae su propósito y un ejemplo mínimo. **No es texto para copiar tal cual**: es la forma que la spec debe respetar.

Las specs se escriben en español y usan los términos de [`Docs/producto/glosario.md`](../../../Docs/producto/glosario.md).

---

## Encabezado

Toda spec empieza con metadatos en un blockquote simple (sin tablas):

```markdown
# SPEC NN — Título corto y descriptivo

> **Estado:** Borrador
> **Tipo:** feat
> **Depende de:** SPEC 01
> **Historias:** HU-08, HU-09
> **Capacidades:** CAP-20
> **Fecha:** AAAA-MM-DD
> **Objetivo:** Una sola oración. Si necesitas dos, la funcionalidad es demasiado grande.
```

**Estados válidos:** `Borrador`, `En revisión`, `Aprobado`, `Implementado`, `Obsoleto`.

- Solo una persona cambia el estado a `Aprobado`. Nunca el agente.
- `/spec-impl` cambia a `Implementado` cuando todos los criterios pasan.

**Tipo:** prefijo de Conventional Commits que tendrán la rama y el PR: `feat`, `fix`, `refactor`, `docs`, `chore`, `ci`, `security`. La rama será `<tipo>/spec-NN-slug`.

**Historias y Capacidades:** IDs de [`Docs/producto/alcance.md`](../../../Docs/producto/alcance.md). Si la spec no cubre ninguna (p. ej. una tarea técnica), escribe `—` y explica por qué en la sección 1.

**Objetivo:** una oración que una persona lee en 5 segundos y entiende qué se va a construir.

---

## Sección 1 — Por qué existe esta spec (opcional)

Para specs que toman decisiones no obvias o rompen patrones del proyecto: el **por qué** del trabajo, no el qué.

En specs simples, se omite.

---

## Sección 2 — Alcance

Dos sub-bloques explícitos. **Ambos son obligatorios.**

```markdown
## Alcance

**Entra:**

- Registrar un gasto con monto, categoría, cuenta y fecha.
- Rechazar fechas futuras.

**Fuera de alcance (para specs futuras):**

- Movimientos recurrentes (v1.3).
- Adjuntar recibos.
```

**Por qué importa "fuera":** guarda lo que se mencionó en las preguntas y se decidió dejar para después. Sin ese registro, durante la implementación aparece la tentación de meterlo "de paso".

---

## Sección 3 — Modelo de datos

Las estructuras concretas que aparecen o cambian, por capa. Código real y corto, no pseudocódigo.

````markdown
## Modelo de datos

**Dominio** — `internal/movements/domain/movement.go`

```go
type Movement struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	AccountID uuid.UUID
	Kind      Kind  // KindIncome | KindExpense
	Amount    Money // entero en la unidad mínima + código ISO 4217
	Date      time.Time
}
```

**Migración** — `db/migrations/00004_create_movements.sql` (nueva; nunca editar una ya mergeada).

Convenciones (ejemplo ilustrativo; las reales las decide cada spec):

- Montos como enteros en la unidad mínima de la moneda; nunca `float`.
- Toda tabla de datos del usuario lleva `user_id` con índice.
````

Si la funcionalidad no introduce datos nuevos, dilo: _"Esta spec no introduce estructuras de datos nuevas. Reutiliza el modelo de SPEC 01."_

---

## Sección 4 — Impacto en seguridad

**Obligatoria (SEC-1). No puede quedar vacía.** Si de verdad no hay impacto, explica por qué en una oración.

```markdown
## Impacto en seguridad

- **Datos que toca:** movimientos (Confidencial) y su descripción libre (posiblemente sensible). Ver SEC-0 §1.
- **Quién puede acceder:** solo el dueño del movimiento. Un ID ajeno responde 404.
- **Escenarios de SEC-0:** A2 (acceso a datos ajenos), A11 (redondeo), A12 (cambios sin rastro).
- **Controles:**
  - `application` verifica `in.Actor` contra `UserID` en cada caso de uso que recibe un ID.
  - Monto validado en el dominio: positivo, moneda igual a la de la cuenta.
- **Datos nuevos que se piden:** ninguno fuera de la lista permitida de SEC-0 §2.
```

Cada control listado aquí debe tener al menos un paso del plan y un criterio de aceptación que lo verifique.

---

## Sección 5 — Plan de implementación

Pasos numerados, agrupados por capa de la arquitectura hexagonal en este orden:

1. **Migración** (si hay cambio de esquema) — `db/migrations/000NN_*.sql` con bloque `-- +goose Down`.
2. **Dominio** — entidad, invariantes, sentinelas de error. Pruebas de dominio.
3. **Puertos** — `port/in` (servicio, commands, queries) y `port/out` (repositorios, dependencias externas).
4. **Aplicación** — casos de uso con autorización por `in.Actor`. Pruebas con stubs de `port/out`.
5. **Adaptador de salida** — `adapter/out/postgres`, con su modelo GORM y mappers. Errores de GORM traducidos a sentinelas.
6. **Adaptador de entrada** — handlers, DTOs, rutas con `guard(...)`, traducción de errores. Pruebas de handler.
7. **Composition root** — `<modulo>_module.go` y registro en `cmd/finanzas/main.go`.

```markdown
## Plan de implementación

1. **Dominio** — Crear `internal/movements/domain/movement.go` con `NewMovement` y sus invariantes. Pruebas: monto ≤ 0 y fecha futura devuelven `ValidationError`.
2. **Puertos** — Definir `in.MovementService` con `Create(ctx, actor, cmd)` y `out.MovementRepository`.
3. ...
```

**Reglas:**

- Cada paso deja el sistema compilando y con `make ci` en verde. Debe poder commitearse solo.
- Cada paso incluye sus pruebas. No hay un paso final de "agregar pruebas".
- Si un paso pasa de 30–50 líneas de código, divídelo.
- Omite las capas que la spec no toca; no las rellenes.
- El último paso no es "probar todo": eso son los criterios de aceptación.

---

## Sección 6 — Criterios de aceptación

Checklist booleano. Cada ítem se verifica con sí o no, y dice **cómo** se verifica.

```markdown
## Criterios de aceptación

- [ ] `make ci` pasa.
- [ ] `POST /api/v1/movements` con fecha de mañana responde 422 con `fields.date`. (prueba de handler)
- [ ] Un usuario que pide el movimiento de otro recibe 404. (prueba de aplicación)
- [ ] `make migrate-up` y `make migrate-reset` funcionan con la migración nueva. (job Migraciones)
```

**Siempre incluir:**

- `make ci` pasa.
- Un criterio por cada control de la sección de seguridad.
- Si hay migración: aplica y revierte limpio.

**Anti-patrones:** "que funcione bien", "buena UX", "sin bugs".

---

## Sección 7 — Decisiones tomadas y descartadas

La sección más valiosa dentro de tres meses. Guarda **lo que se consideró**, no solo lo elegido.

```markdown
## Decisiones

- **Sí:** montos como `int64` en la unidad mínima. Evita errores de redondeo (A11).
- **No:** `decimal` de una librería externa. Dependencia nueva sin necesidad para el MVP.
- **Sí:** anular en vez de borrar. Glosario: "Los movimientos no se borran; se anulan".
```

Cada decisión con una razón breve. Las decisiones sin razón son las primeras que se cuestionan.

---

## Sección 8 — Riesgos identificados (opcional)

Solo si hay riesgos no obvios.

```markdown
## Riesgos

| Riesgo | Mitigación |
| --- | --- |
| La fuente de tasas no responde | Se usa la última tasa conocida y se marca como desactualizada (A18). |
```

---

## Sección final — Lo que NO está en esta spec (refuerzo)

Repite al final lo que **no** se hará. Es deliberado: sirve de recordatorio a quien lee solo las últimas líneas.

```markdown
## Lo que **no** está en esta spec

- Movimientos recurrentes.
- Importación de extractos.

Cada uno, si llega, va en su propia spec.
```

---

## Reglas globales del documento

- **Una oración por idea.** Si una oración tiene dos comas y un punto y coma, divídela.
- **Nombres concretos.** No "el módulo de movimientos": `internal/movements/`. No "una ruta": `POST /api/v1/movements`.
- **Sin TODOs.** Un TODO significa que la decisión no se tomó. Tómala o anótala como decisión pendiente con su razón.
- **Sin código ejecutable largo.** Fragmentos cortos para ilustrar estructuras sí; funciones completas no.
- **Términos del glosario.** Usa las palabras de `Docs/producto/glosario.md` con su significado. Si la spec necesita un término que no está, defínelo en la propia spec.
- **Sin datos reales ni secretos** en ejemplos.
- **Markdown estándar** que se vea bien en GitHub.
