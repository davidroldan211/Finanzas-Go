# Política de seguridad

## Reportar una vulnerabilidad

**No abras un issue público.** Reporta en privado desde
[Security → Report a vulnerability](https://github.com/davidroldan211/Finanzas-Go/security/advisories/new).

Incluye, si puedes:

- Ruta o componente afectado.
- Pasos para reproducir (sin datos reales de terceros).
- Impacto que observaste.

## Qué esperar

| Paso | Plazo objetivo |
|---|---|
| Acuse de recibo | 3 días hábiles |
| Evaluación inicial y severidad | 7 días hábiles |
| Corrección de severidad crítica o alta | 30 días |

El reporte se mantiene privado hasta que haya una corrección publicada. Si quieres, se te reconoce en el advisory.

## Versiones soportadas

El proyecto está en desarrollo previo a la versión 1.0. Solo se corrige la rama `main`.

## Prácticas del proyecto

- Escaneo de secretos con bloqueo en push, CodeQL, `govulncheck` y revisión de dependencias en cada PR.
- Acciones de GitHub fijadas por SHA y actualizadas por Dependabot.
- Puntos de control de seguridad por fase: ver [Docs/roadmap.html](Docs/roadmap.html).
