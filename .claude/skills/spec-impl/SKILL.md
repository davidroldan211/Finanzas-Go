---
name: spec-impl
description: Implementa una spec aprobada de Finanzas-Go. Valida que el estado sea "Aprobado", crea la rama <tipo>/spec-NN-slug desde main, implementa paso a paso por capas con pausas para revisar el diff, verifica los criterios con make ci, marca la spec como Implementado y abre el PR. Nunca hace merge.
disable-model-invocation: true
argument-hint: <NN-nombre-de-spec>
allowed-tools: Read, Glob, Grep, Edit, Write, AskUserQuestion, Bash(git status:*), Bash(git branch:*), Bash(git switch:*), Bash(git checkout:*), Bash(git fetch:*), Bash(git pull:*), Bash(git log:*), Bash(git diff:*), Bash(git add:*), Bash(git commit:*), Bash(git show:*), Bash(git push:*), Bash(make:*), Bash(go build:*), Bash(go test:*), Bash(go vet:*), Bash(gofmt:*), Bash(docker compose ps:*), Bash(gh pr create:*), Bash(gh pr view:*), Bash(gh pr checks:*), Bash(cat:*), Bash(ls:*)
---

# /spec-impl — Implementador de specs aprobadas

## Contexto de la sesión

Estado del repositorio:
!`git status --short`

Rama actual:
!`git branch --show-current`

Specs disponibles:
!`ls specs/ 2>/dev/null || echo "La carpeta specs/ no existe"`

Configuración del flujo:
!`cat specs/.spec-config.yml 2>/dev/null || echo "Sin archivo de configuración: AutoCreateBranch: true, AutoCommitSteps: false"`

---

## Instrucciones

Sigue estas cinco fases en orden estricto. **No avances a la siguiente si la anterior no terminó bien.**

---

### Fase 1 — Identificar la spec

El argumento recibido es: `$ARGUMENTS`

Si `$ARGUMENTS` está vacío:

- Muestra las specs de `specs/` (las tienes arriba).
- Pide el nombre exacto de la spec.
- Para y espera la respuesta.

Si `$ARGUMENTS` tiene valor:

- Busca el archivo en `specs/`. La persona puede haber escrito el nombre completo (`03-registrar-gasto`), solo el número (`03`) o solo el slug (`registrar-gasto`). Encuentra el archivo en cualquiera de esos casos.
- Si no lo encuentras, muestra las specs disponibles y pide corregir el nombre.
- Si lo encuentras, sigue a la Fase 2.

---

### Fase 2 — Validar el estado

Lee la spec. Busca la línea `**Estado:**` del encabezado.

**Regla absoluta:** solo continúas si el estado es `Aprobado` (o su equivalente claro en otro idioma, como `Approved`).

| Estado | Acción |
| --- | --- |
| `Aprobado` | Sigue a la Fase 3. |
| `Borrador`, `En revisión` | Para. Muestra el mensaje de abajo. |
| `Implementado`, `Obsoleto` | Para. Muestra el mensaje de abajo. |
| Sin línea de estado o valor no reconocido | Para. Di que el archivo no sigue el formato de `.claude/skills/spec/template.md`. |

Si dudas de si un valor significa "aprobado", **no supongas**: para y pide aclararlo.

Además, verifica que la spec tenga **`**Tipo:**`** y una sección **Impacto en seguridad** no vacía. Si falta alguna, para y pide completarla con `/spec`.

**Mensaje cuando el estado no es Aprobado:**

```
❌ No puedo implementar esta spec.

Estado actual: [ESTADO ENCONTRADO]
Solo implemento specs en estado "Aprobado".

Para continuar:
  1. Si la spec está lista, ábrela y cambia el estado a "Aprobado".
     Ese cambio lo hace una persona, no el agente.
  2. Si todavía necesita trabajo, retómala con /spec [nombre].
```

No ofrezcas alternativas ni "puedo empezar igual si quieres". El bloqueo es intencional.

---

### Fase 3 — Crear la rama y cambiarse a ella

1. **Revisa el árbol de trabajo.** Si `git status --short` del contexto **no está vacío**, para, muestra los cambios pendientes y pregunta:

   ```
   ⚠️ Hay cambios sin commitear.
   Cambiar de rama los llevaría consigo. ¿Qué hacemos?
     1. Los commiteas o guardas tú y vuelves a ejecutar el comando  (recomendado)
     2. Seguir igual: los cambios viajan a la rama nueva
   ```

   Espera la respuesta. **No hagas stash ni commit por la persona** salvo que lo pida.

2. **Nombre de la rama:** `<tipo>/spec-NN-slug`, con `<tipo>` del campo `**Tipo:**` de la spec y `NN-slug` del nombre del archivo sin extensión.
   - `03-registrar-gasto.md` con `Tipo: feat` → `feat/spec-03-registrar-gasto`
   - `04-tope-intentos-login.md` con `Tipo: security` → `security/spec-04-tope-intentos-login`

3. Lee `AutoCreateBranch` de la configuración. Si falta o el valor no se reconoce, es `true`. Solo un `false` explícito lo desactiva.

   **Si es `true`:**
   - Si la rama **no existe**: actualiza `main` y crea la rama desde ahí:
     ```bash
     git switch main
     git pull --ff-only
     git switch -c <tipo>/spec-NN-slug
     ```
   - Si **ya existe**: se está retomando trabajo. Cámbiate a ella, lee `git log --oneline main..HEAD` y di qué pasos del plan parecen hechos y desde cuál propones retomar. Los commits `spec-NN-step-MM` marcan los pasos terminados. Espera confirmación antes de implementar.

   **Si es `false`:** pregunta antes de tocar git:
   ```
   AutoCreateBranch está en false.
   ¿Creo la rama <tipo>/spec-NN-slug desde main y me cambio a ella? [s/N]
   ```
   Si responde no, **no crees ninguna rama**. Recuerda que en `main` no se puede hacer push directo y pide confirmación explícita para seguir en la rama actual.

4. Confirma:

   ```
   ✅ Lista para implementar.

   Spec:   specs/NN-slug.md
   Rama:   <tipo>/spec-NN-slug  (activa)
   Estado: Aprobado
   ```

5. **Todavía no implementes.** Muestra el resumen de la spec:
   - **Objetivo**, **Historias** y **Capacidades**.
   - **Alcance** (entra / fuera).
   - **Impacto en seguridad**: controles que hay que implementar.
   - **Plan de implementación**.
   - **Criterios de aceptación**.

---

### Fase 4 — Implementar paso a paso

Después del resumen, di:

```
Voy a implementar la spec siguiendo el plan al pie de la letra.
Pararé después de cada paso para que revises el diff.
Cada paso se commiteará como spec-NN-step-MM.   (← solo si AutoCommitSteps es true)

¿Empiezo con el Paso 1?
```

Espera confirmación explícita ("sí", "dale", "adelante" o equivalente).

**Reglas del proyecto durante toda la implementación:**

- **Arquitectura hexagonal** (ver `CLAUDE.md`): `domain`, `port` y `application` nunca importan `gin`, `gorm`, `config` ni `httpx`.
- **Autorización en `application`** con `in.Actor`: todo caso de uso que recibe un ID de recurso verifica que pertenece al actor. Un recurso ajeno responde igual que uno inexistente (404).
- **Migraciones:** archivo nuevo con `make migrate-create name=...` y bloque `-- +goose Down`. **Nunca edites una migración ya mergeada.**
- **Errores de GORM** se traducen a sentinelas de dominio dentro del adaptador de salida.
- **Pruebas en el mismo paso** que el código que prueban. Sin datos reales ni secretos en pruebas.
- **Términos del glosario** (`Docs/producto/glosario.md`) en nombres y mensajes.

**Verificación de cada paso** antes de pausar:

```bash
go build ./...
go test ./internal/<modulo>/...
make arch
```

Si algo falla, corrígelo dentro del mismo paso. No pases al siguiente con el árbol en rojo.

**Commits:** controlados por `AutoCommitSteps`. Si falta o no se reconoce, es `false`.

- **`false` (por defecto):** no commitees durante los pasos. El commit único se hace en la Fase 5.
- **`true`:** commitea cada paso al terminarlo, antes de pausar:
  - Agrega solo los archivos del paso: `git add -- <archivo> <archivo> …`. Nunca `git add -A` ni `git add .`.
  - Mensaje `spec-NN-step-MM`, sin cuerpo (dos dígitos cada número). El squash merge del PR lo reemplaza en `main`.
  - Si la persona pide ajustes al paso en revisión, inclúyelos en el mismo commit (`git add -- <archivos>` + `git commit --amend --no-edit`). Solo enmiendas el commit del paso en revisión.
  - Si un hook rechaza el commit, para y muestra el error. Nunca uses `--no-verify`.

**Ritmo:**

- Implementa un paso.
- Corre la verificación del paso.
- Si `AutoCommitSteps` es `true`, commitea.
- Muestra qué archivos tocaste y qué hiciste. Si commiteaste: `Commit spec-NN-step-MM (<hash corto>)`.
- Di: `Paso N terminado. ¿Revisas el diff y sigo con el Paso N+1?`
- Espera confirmación.

**La regla principal:** implementa lo que dice la spec. Si algo te parece mejorable, menciónalo como observación pero implementa lo acordado. Los cambios a la spec van en la spec, no en el código por sorpresa.

**Si encuentras una ambigüedad** que la spec no resuelve: para, descríbela exactamente, presenta 2 o 3 opciones concretas y espera la decisión. No improvises.

**Si la persona pide algo fuera del alcance:** recuérdale que está fuera de esta spec, sugiere anotarlo para una spec futura y no lo implementes en esta rama.

---

### Fase 5 — Verificar, cerrar y abrir el PR

Cuando termines el último paso:

1. **Pipeline local:**
   ```bash
   make ci
   ```
   Si la spec agrega una migración y la base local está arriba (`docker compose ps`), corre también `make migrate-up`, `make migrate-reset` y `make migrate-up`. Si la base no está arriba, dilo: el job **Migraciones** del PR lo verificará.

2. **Criterios de aceptación uno por uno.** Para cada criterio, di cómo lo verificaste (prueba, comando, salida) y si pasa. Si alguno no pasa, **para**: muéstralo y pregunta cómo seguir. No marques la spec como implementada.

3. **Actualiza la spec:**
   - `**Estado:** Implementado`.
   - Marca `[x]` los criterios verificados.

4. **Commit:**
   - Con `AutoCommitSteps: false`: un solo commit con todos los archivos de la spec, agregados uno por uno (`git add -- …`). Mensaje Conventional Commits: `<tipo>(<modulo>): <objetivo en minúscula>`, y en el cuerpo `Spec: specs/NN-slug.md` y las historias.
   - Con `AutoCommitSteps: true`: un commit final `spec-NN-cierre` con la spec actualizada.

5. **Push y PR:**
   ```bash
   git push -u origin <tipo>/spec-NN-slug
   ```
   Abre el PR con `gh pr create --base main`:
   - **Título:** Conventional Commits, `<tipo>(<modulo>): <objetivo>`; es el mensaje que quedará en `main`.
   - **Cuerpo**, siguiendo `.github/pull_request_template.md`:
     - Enlace a `specs/NN-slug.md`, historias y capacidades. `Closes #N` si la spec referencia un issue.
     - Qué cambia, por capa.
     - Cómo se probó: la evidencia de cada criterio del paso 2.
     - Checklist de seguridad con los controles de la spec.
   - **Etiqueta:** `historia` si cubre HU; `seguridad` si el tipo es `security` o la spec atiende un escenario crítico de SEC-0.

6. Informa:

   ```
   ✅ Spec implementada y PR abierto.

   PR:      <url>
   Rama:    <tipo>/spec-NN-slug
   Estado:  Implementado

   Los checks del PR están corriendo. El merge lo haces tú.
   ```

---

## Reglas duras

- **Nunca hagas merge** del PR, ni con `gh pr merge` ni de ninguna otra forma. El merge siempre lo hace la persona.
- **Nunca hagas push a `main`** ni `git push --force`.
- **Nunca cambies una spec a `Aprobado`.** Solo a `Implementado`, y solo en la Fase 5 con todos los criterios verificados.
- **Nunca edites una migración ya mergeada.**
- **Nunca uses `--no-verify`** ni desactives controles del pipeline para que algo pase.
- **No implementes nada fuera del alcance de la spec.**

---

## Resumen del comportamiento esperado

```
/spec-impl 03-registrar-gasto   (Estado: Aprobado, Tipo: feat)

  Fase 1  →  Encuentra specs/03-registrar-gasto.md
  Fase 2  →  Estado "Aprobado", Tipo y seguridad presentes → ✅ sigue
  Fase 3  →  main actualizado → git switch -c feat/spec-03-registrar-gasto
             Muestra objetivo, alcance, seguridad, plan y criterios
  Fase 4  →  Implementa paso a paso por capas, con pausas y verificación
  Fase 5  →  make ci, criterios uno por uno, Estado: Implementado,
             commit, push, PR. No hace merge.

/spec-impl 04-tope-intentos-login   (Estado: Borrador)

  Fase 1  →  Encuentra specs/04-tope-intentos-login.md
  Fase 2  →  Estado "Borrador" → ❌ para y muestra el mensaje
             No crea rama, no toca código
```

**`AutoCreateBranch`** (`specs/.spec-config.yml`, por defecto `true`): crea la rama sin preguntar. Con `false`, pide `[s/N]`.

**`AutoCommitSteps`** (mismo archivo, por defecto `false`): un solo commit al final. Con `true`, un commit por paso `spec-NN-step-MM`.
