---
name: spec
description: Diseña specs de Finanzas-Go con el método spec-driven. Hace preguntas antes de proponer estructura y escribe la spec en specs/ con la plantilla del proyecto (trazabilidad HU/CAP, impacto en seguridad, plan por capas hexagonales). Úsala antes de escribir código de una funcionalidad.
disable-model-invocation: true
argument-hint: 'descripción corta de la funcionalidad o requisito'
allowed-tools: Read, Glob, Grep, Write, AskUserQuestion, Bash(ls:*), Bash(cat:*), Bash(date:*)
---

# /spec — Diseñador guiado de specs

## Contexto de la sesión

Fecha de hoy (úsala en el encabezado, nunca la adivines):
!`date +%F`

Specs que ya existen:
!`ls specs/ 2>/dev/null || echo "La carpeta specs/ todavía no existe"`

---

Esta skill ayuda a producir una spec útil con el método spec-driven. **Aquí no se escribe código.** Tu trabajo es ayudar a la persona a aclarar qué quiere construir, preguntar cuando algo no está suficientemente definido y escribir la spec hasta que esté lista para guardarse en `specs/`.

## Filosofía

Una spec no es documentación decorativa: es el contrato que guía la implementación. Si la spec es vaga, el código improvisa. Por eso este flujo es **lento a propósito al definir** y **rápido al escribir**.

Lee `template.md` (en la misma carpeta que esta skill) para ver la estructura completa. Apóyate en ella en cada paso.

## Flujo

- Sigue las cuatro fases en orden. **Nunca te saltes la Fase 2**: las preguntas son el punto. Si la persona quiere ir más rápido, recuérdale que el costo de una mala spec se paga después en código.
- Responde en el idioma del prompt inicial. El proyecto trabaja en español.

### Fase 1 — Entender el contexto

Antes de preguntar por la funcionalidad, carga el contexto del proyecto:

1. `CLAUDE.md`: arquitectura hexagonal, reglas de migraciones, flujo de ramas y PRs.
2. `Docs/producto/vision.md`: para quién es la app y qué no es.
3. `Docs/producto/alcance.md`: capacidades (CAP-xx), historias (HU-xx) y las decisiones de la sección 5. Identifica qué HU y CAP cubre la funcionalidad. Si no corresponde a ninguna, pregúntalo en la Fase 2.
4. `Docs/producto/glosario.md`: usa estos términos con este significado.
5. `Docs/seguridad/sec-0-datos-y-amenazas.md`: clasificación de datos, datos que nunca se piden y escenarios de amenaza (A1–A18).
6. Revisa el listado de `specs/` del contexto. Si hay specs anteriores, lee al menos las dos más recientes para tomar sus convenciones.
7. Mira el código de los módulos que la funcionalidad va a tocar o imitar (`internal/<modulo>/`).

Si `$ARGUMENTS` viene vacío, pide una descripción **de una sola oración**. Si no cabe en una oración, es la primera señal de que es demasiado grande: propone dividirla.

### Fase 2 — Aclarar con preguntas

La fase más importante. Tu trabajo aquí es **detectar ambigüedades y preguntar**, no suponer.

Pregunta en bloques de 3 a 5 a la vez, con `AskUserQuestion`: 2–4 opciones por pregunta, tu recomendación primero y marcada como "(Recomendado)". Después de cada bloque, espera la respuesta.

**Categorías que siempre debes considerar:**

- **Alcance:** qué entra y qué NO. Qué partes se dejan para otra spec.
- **Trazabilidad:** qué HU y CAP de `alcance.md` cubre. Si cubre varias historias, ¿conviene una spec por historia?
- **Datos:** qué estructuras nuevas aparecen, cómo se llaman (según el glosario), en qué módulo viven.
- **Persistencia:** ¿hay cambio de esquema? ¿Qué migración nueva? ¿Cómo se revierte?
- **Integración:** ¿depende de specs anteriores? ¿Modifica algo existente o solo agrega? ¿Necesita leer datos de otro módulo? (Patrón: adaptador de salida propio, nunca importar la infraestructura ajena.)
- **API y estados:** rutas, códigos de respuesta, qué pasa cuando falla.
- **Seguridad:** quién puede acceder a cada recurso; qué datos de SEC-0 toca y de qué nivel; qué escenarios A1–A18 aplican; si pide algún dato de la lista de "nunca se piden".
- **Riesgos:** qué puede romper esto; qué pasa en el caso degradado.
- **Decisiones cerradas:** algo que la persona ya decidió y no quiere reabrir. Las decisiones de `alcance.md` §5 están cerradas: no las vuelvas a preguntar.

**Cómo formular las preguntas:**

- Concretas, no abiertas. ❌ "¿Cómo imaginas la persistencia?" → ✅ "¿El monto se guarda como `int64` en la unidad mínima o como `NUMERIC(18,2)`?"
- Si una respuesta abre una caja de Pandora (p. ej. "y también importar extractos"), señala que merece su propia spec y pregunta si queda fuera de esta.

**Cuándo dejar de preguntar:**

Cuando puedas responder estas cuatro preguntas sin suponer nada:

1. ¿Qué archivos aparecen o cambian, y en qué capa?
2. ¿Cuál es el primer paso ejecutable y cuál el último?
3. ¿Cómo verifico que la funcionalidad está terminada?
4. ¿Quién puede acceder a cada dato nuevo, y qué control lo garantiza?

Si todavía no puedes responder alguna, sigue preguntando.

### Fase 3 — Escribir la spec

Con la Fase 2 cerrada, decide cómo escribir:

**Si ya tienes toda la información** —puedes responder las cuatro preguntas sin suponer— **no vayas sección por sección.** Escribe la spec completa y pasa directo a la Fase 4. No pidas confirmación por sección ni muestres un borrador para aprobar: la persona ya respondió todo y revisa el archivo guardado.

**Solo si falta información** (la Fase 2 se cortó, una respuesta fue vaga o una sección no se puede escribir sin inventar), desarrolla las secciones **una por una**, mostrando cada una y esperando confirmación.

En ambos casos, el orden es:

1. **Encabezado:** estado, tipo, dependencias, historias, capacidades, fecha y objetivo de una oración.
2. **Alcance:** qué entra y qué NO, explícito.
3. **Modelo de datos:** estructuras concretas por capa y migración. Si no hay datos nuevos, dilo.
4. **Impacto en seguridad:** obligatoria; nunca vacía.
5. **Plan de implementación:** pasos numerados por capa, cada uno con sus pruebas y dejando `make ci` en verde.
6. **Criterios de aceptación:** checklist booleano, con `make ci` y un criterio por cada control de seguridad.
7. **Decisiones tomadas y descartadas:** con su razón.
8. **Riesgos:** solo si aplica.
9. **Lo que NO está en esta spec:** refuerzo final.

**Errores comunes a evitar:**

- Criterios de aceptación no verificables ("que funcione bien").
- Pasos del plan que no están en el alcance.
- Suponer nombres de archivos o estructuras que la persona no confirmó.
- Saltarse la sección de decisiones o la de seguridad.
- Un control de seguridad sin prueba que lo verifique.

### Fase 4 — Guardar la spec

1. Calcula el siguiente número a partir del listado de `specs/`: el mayor existente más uno, con dos dígitos. Si `specs/` está vacía o no existe, empieza en `01-`.
2. Genera un slug corto en kebab-case a partir del objetivo (p. ej. `registrar-gasto`).
3. Usa la fecha del contexto de la sesión para `**Fecha:**`. **Nunca escribas una fecha que no leíste de ahí.**
4. Escribe el archivo directamente en `specs/NN-slug.md`. **No pidas permiso para escribirlo ni confirmación del nombre**: anuncia la ruta al final. Solo pregunta si el archivo ya existe.
5. Marca el estado como `Borrador`. **Nunca lo marques como `Aprobado`**: eso lo hace la persona después de releerla.
6. Si el encabezado lista dependencias (`**Depende de:** SPEC 01`), verifica que cada spec exista en `specs/`. Si no existe, dilo en vez de dejar una referencia rota.
7. Verifica que cada HU y CAP del encabezado exista en `Docs/producto/alcance.md`.
8. **Crea el archivo de configuración si no existe.** Si falta `specs/.spec-config.yml`, créalo con el contenido de abajo. Si ya existe, **no lo toques**.

   ```yaml
   # Configuración del flujo de specs
   #
   # AutoCreateBranch — si /spec-impl crea la rama automáticamente.
   #   true  (por defecto) → crea <tipo>/spec-NN-slug desde main sin preguntar
   #   false               → pide confirmación [s/N] antes de crearla
   AutoCreateBranch: true

   # AutoCommitSteps — si /spec-impl commitea cada paso terminado.
   #   false (por defecto) → un solo commit al final, con mensaje Conventional Commits
   #   true                → un commit por paso: spec-NN-step-MM
   AutoCommitSteps: false
   ```

9. Confirma a la persona:
   - Ruta del archivo creado.
   - Recordatorio: la spec está en `Borrador`. Cámbiala a `Aprobado` cuando la hayas releído.
   - Si acabas de crear `specs/.spec-config.yml`, menciona sus dos opciones y sus valores por defecto.
   - Siguiente paso: cuando esté aprobada, ejecutar `/spec-impl NN-slug`.
   - **Para aquí.** No propongas implementarla ni hagas nada más.

## Reglas duras

- **Nunca escribas código en este comando.** Solo el `.md` de la spec y, si falta, `specs/.spec-config.yml`.
- **Nunca propongas implementar la spec después de guardarla.**
- **Nunca supongas decisiones que la persona no confirmó.** Si falta información, pregunta en la Fase 2.
- **No vuelvas a preguntar en la Fase 3 lo que ya se respondió en la Fase 2.**
- **La sección de impacto en seguridad nunca queda vacía.**
- **No pidas datos de la lista "nunca se piden" de SEC-0 §2.** Si la funcionalidad parece necesitarlos, dilo y para.
- **Sin datos reales ni secretos** en ejemplos de la spec.
- **Si la persona quiere saltarse la Fase 2**, recuérdale: "Las preguntas de ahora ahorran horas después. ¿Seguro que quieres saltarlas?". Si insiste, respétalo y regístralo en las decisiones ("Definición rápida sin aclaración detallada").
- **Si la funcionalidad es demasiado grande** (no cabe en una oración, toca más de tres módulos o pide decisiones en cuatro o más dominios), propone dividirla en varias specs antes de seguir.

## Tono al preguntar

Directo y específico. No te disculpes por preguntar. Nada de "si no te molesta…". La persona invocó esta skill porque quiere preguntas.

Ejemplo de un bloque bien formado:

> Antes de escribir el modelo de datos necesito aclarar tres cosas:
>
> 1. **Monto.** ¿`int64` en la unidad mínima o `NUMERIC(18,2)`? Recomiendo `int64`: evita errores de redondeo (A11).
> 2. **Anulación.** ¿Columna `voided_at` o tabla de historial aparte? Recomiendo tabla de historial: cubre también las correcciones (A12).
> 3. **Acceso.** ¿Un admin puede ver movimientos de otros usuarios? La visión dice que no; confirmo que no.

## Argumentos

`$ARGUMENTS` es **la descripción de la funcionalidad**, no el nombre del archivo. Úsala como punto de partida de la Fase 1 y deriva el slug del objetivo en la Fase 4.

Excepción: si `$ARGUMENTS` es un solo token en kebab-case sin espacios (p. ej. `/spec registrar-gasto`), úsalo como slug **y** como semilla de la descripción, sin pedir confirmación.

Si se invocó `/spec` sin argumentos, empieza pidiendo la descripción de una oración.
