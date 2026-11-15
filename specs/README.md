# Specs

Cada funcionalidad de Finanzas-Go nace como una spec en esta carpeta antes de escribir código. La spec es el contrato: dice qué se construye, qué no, cómo se protege y cómo se verifica.

## Ciclo

| Etapa | Quién | Cómo | Resultado |
|---|---|---|---|
| 1. Idea | Persona | Historia en [`Docs/producto/alcance.md`](../Docs/producto/alcance.md) o issue | HU-xx / CAP-xx |
| 2. Spec | IA + persona | `/spec <descripción>` | `specs/NN-slug.md` en **Borrador** |
| 3. Revisión | **Persona** | Relee la spec y cambia el estado a **Aprobado** | Spec lista |
| 4. Implementación | IA + persona | `/spec-impl NN-slug`, paso a paso con pausas | Rama `<tipo>/spec-NN-slug` |
| 5. Cierre | IA | `make ci`, criterios uno por uno, estado **Implementado** | PR abierto |
| 6. Merge | **Persona** | Revisa el PR con los checks en verde y hace squash merge | Cambio en `main` |

Las etapas 3 y 6 son humanas: la IA nunca aprueba una spec ni hace merge.

Plantilla: [`.claude/skills/spec/template.md`](../.claude/skills/spec/template.md).

## Estados

`Borrador` → `En revisión` → `Aprobado` → `Implementado`. `Obsoleto` para specs que ya no aplican.

## Listo para implementar

Una spec está lista cuando:

- [ ] Su estado es **Aprobado**, puesto por una persona.
- [ ] Tiene **Tipo**, **Historias** y **Capacidades** en el encabezado.
- [ ] El objetivo cabe en una oración.
- [ ] El alcance dice qué entra y qué no.
- [ ] La sección **Impacto en seguridad** no está vacía y cada control tiene un criterio de aceptación.
- [ ] Cada criterio de aceptación se verifica con sí o no.

`/spec-impl` se niega a empezar si falta el estado, el tipo o la sección de seguridad.

## Terminado

Una spec está terminada cuando:

- [ ] `make ci` pasa en local y el check **CI OK** pasa en el PR.
- [ ] Todos los criterios de aceptación están marcados `[x]` con su evidencia en el PR.
- [ ] El estado de la spec es **Implementado**.
- [ ] El PR está mergeado a `main`.
- [ ] Si la spec cambió algo estructural, `CLAUDE.md` lo refleja.

## Configuración

[`.spec-config.yml`](.spec-config.yml):

- `AutoCreateBranch` (por defecto `true`): `/spec-impl` crea la rama sin preguntar.
- `AutoCommitSteps` (por defecto `false`): un solo commit al final; con `true`, un commit por paso.
