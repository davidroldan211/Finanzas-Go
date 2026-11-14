## Qué cambia y por qué

<!-- Una o dos frases. El título del PR sigue Conventional Commits: tipo(alcance): asunto -->

## Spec / issue

<!-- Closes #N · HU-xx · CAP-xx -->

## Cómo se probó

<!-- Pruebas agregadas, comandos ejecutados, casos manuales. -->

## Checklist

- [ ] `make ci` pasa en local
- [ ] Pruebas nuevas o actualizadas para el comportamiento cambiado
- [ ] Si hay cambio de esquema: migración nueva (nunca editar una ya mergeada) con bloque `Down`
- [ ] Documentación actualizada (README, `CLAUDE.md`, `Docs/`) si aplica
- [ ] Términos según el [glosario](../Docs/producto/glosario.md)

## Seguridad

- [ ] Toda ruta nueva tiene guard y verifica que el recurso pertenece al usuario
- [ ] Entradas validadas; errores 5xx no exponen detalles internos
- [ ] Sin secretos, credenciales ni datos reales en código, pruebas o logs
- [ ] Código generado con IA revisado línea por línea
- [ ] Dependencias nuevas justificadas y revisadas
