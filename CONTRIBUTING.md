# Cómo contribuir

## Flujo de trabajo

1. **Todo empieza en un issue** (plantilla *Historia de usuario* o *Error*). Las historias referencian su capacidad del [backlog](Docs/producto/alcance.md).
2. **Funcionalidades: spec primero.** Con Claude Code, `/spec <descripción>` crea `specs/NN-slug.md` en Borrador; una persona la aprueba y `/spec-impl NN-slug` la implementa y abre el PR. Ciclo, "Listo" y "Terminado" en [specs/README.md](specs/README.md). Correcciones pequeñas y cambios de documentación pueden ir sin spec.
3. **Rama desde `main`**, con prefijo según el tipo de cambio:

   | Prefijo | Uso |
   |---|---|
   | `feat/` | Funcionalidad nueva |
   | `fix/` | Corrección de un error |
   | `docs/` | Solo documentación |
   | `refactor/` | Cambio interno sin cambiar comportamiento |
   | `test/` | Solo pruebas |
   | `chore/`, `ci/` | Mantenimiento, dependencias, pipeline |
   | `security/` | Corrección o control de seguridad |

   Con spec: `<tipo>/spec-NN-slug` (p. ej. `feat/spec-03-registrar-gasto`). Sin spec: `feat/hu-08-registrar-gasto`.
4. **Antes de abrir el PR:** `make ci` debe pasar en local.
5. **PR hacia `main`** completando la plantilla. El título sigue [Conventional Commits](https://www.conventionalcommits.org/es/v1.0.0/) porque se convierte en el commit de `main` (squash merge):

   ```
   feat(users): permitir registro público con verificación de correo
   fix(auth): rechazar tokens con algoritmo distinto a HS256
   ```
6. **Merge** solo con el check `CI OK` en verde y las conversaciones resueltas. La rama se borra automáticamente.

## Reglas que no se negocian

- `main` no recibe pushes directos.
- Una migración ya mergeada no se edita: los cambios de esquema van en una migración nueva.
- El núcleo (`domain`, `port`, `application`) no importa infraestructura. Lo verifica `make arch`.
- Sin secretos ni datos reales en el repositorio, en pruebas ni en logs.
- Las vulnerabilidades se reportan en privado: ver [SECURITY.md](SECURITY.md).

## Comandos útiles

```bash
make ci          # formato, go.mod, vet, arquitectura, pruebas con cobertura y build
make golangci    # linter completo (requiere golangci-lint v2)
make vuln        # vulnerabilidades conocidas (requiere govulncheck)
make fmt         # formatea el código
```

## Qué corre en cada PR

| Workflow | Checks |
|---|---|
| CI | Lint, Test, Build, Vulnerabilidades, Migraciones y el agregador **CI OK** |
| Seguridad | CodeQL y revisión de dependencias (también semanal sobre `main`) |
| Título del PR | Conventional Commits |
