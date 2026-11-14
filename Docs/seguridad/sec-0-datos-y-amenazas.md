# SEC-0 · Clasificación de datos y modelo de amenazas

> Punto de control de la fase F0 del [roadmap](../roadmap.html). Estado: **borrador para revisión** · 2026-11-14
> Se revisa de nuevo en SEC-6, antes del lanzamiento.

## 1. Clasificación de datos

| Nivel | Significado |
|---|---|
| **Restringido** | Si se filtra, permite tomar control de cuentas o de la plataforma. Nadie lo lee, ni siquiera los administradores. |
| **Confidencial** | Datos personales o financieros del usuario. Solo el propio usuario lo ve. |
| **Interno** | Datos de operación sin información personal identificable. |
| **Público** | Puede publicarse sin daño. |

| Dato | Ejemplos | Nivel | Notas |
|---|---|---|---|
| Credenciales | Contraseña, códigos de verificación, tokens de sesión | Restringido | Nunca se guardan en claro ni aparecen en logs. |
| Secretos de la plataforma | Clave de firma de tokens, credenciales de base de datos | Restringido | Fuera del repositorio. |
| Datos de identificación | Correo, nombre, apellido | Confidencial | El correo identifica al usuario y sirve para recuperar el acceso. |
| Datos financieros | Cuentas, saldos, movimientos, montos, patrimonio | Confidencial | Revelan nivel de ingresos, hábitos y deudas. |
| Texto libre | Descripciones y notas de movimientos, nombres de categorías y cuentas | Confidencial, posiblemente **sensible** | El usuario puede escribir cosas como "cita psiquiatra", "diezmo" o "abogado divorcio", que revelan salud, religión o vida privada. Se trata como dato sensible. |
| Preferencias | Moneda principal, zona horaria | Confidencial | Bajo impacto, pero personal. |
| Eventos de seguridad | Intentos de inicio de sesión, cambios de contraseña | Interno | Contienen IP y correo: se minimizan y tienen retención limitada. |
| Métricas agregadas | Número de usuarios, movimientos por día | Interno | Sin datos individuales. |
| Tasas de cambio | COP/USD por fecha | Público | — |

## 2. Datos que nunca se piden

Decisión de producto. Una spec que necesite alguno de estos datos requiere revisar este documento primero.

- Número completo de cuentas bancarias o de tarjetas.
- Usuario o contraseña de la banca en línea.
- Número de cédula u otro documento de identidad.
- Fecha de nacimiento, dirección, teléfono.
- Ubicación del dispositivo.
- Fotos de recibos o documentos (al menos hasta F8).

Para que el usuario distinga sus cuentas basta con un nombre libre ("Bancolombia ahorros", "Visa ****1234" si él mismo decide escribirlo).

## 3. Normativa a evaluar

> No es asesoría legal. Esta lista identifica qué revisar con un profesional antes de F6.

- **Ley 1581 de 2012 y Decreto 1377 de 2013 (Colombia)**: protección de datos personales. Requieren autorización del titular, finalidad informada, política de tratamiento y canales para consultar, corregir y suprimir datos.
- **Datos sensibles** según esa ley: el texto libre (sección 1) puede contenerlos. Evaluar si se necesita una autorización específica o un aviso.
- **Registro Nacional de Bases de Datos (SIC)**: evaluar si aplica según el tipo y tamaño de la entidad que opere la app.
- **Transferencia internacional de datos**: si la base de datos o algún proveedor (correo, hosting, IA) está fuera de Colombia.
- **Usuarios fuera de Colombia**: el MVP apunta a Colombia, pero cualquiera puede registrarse. Decidir si se limita o si se evalúa otra normativa.

Consecuencia para el alcance: HU-21 (exportar y borrar mis datos) probablemente pase de "Should" a "Must" antes del lanzamiento.

## 4. Modelo de amenazas de alto nivel

### Activos a proteger

1. Datos financieros y personales de cada usuario.
2. Acceso a las cuentas de usuario.
3. Privilegios de administrador.
4. Exactitud de los saldos (un saldo alterado sin que el usuario lo note también es un incidente).
5. Disponibilidad del servicio.
6. Código, pipeline y secretos de la plataforma.

### Actores

| Actor | Motivación | Capacidad |
|---|---|---|
| Atacante externo oportunista | Credenciales reutilizadas, datos para vender | Automatización, listas de contraseñas filtradas |
| Otro usuario registrado | Curiosidad o mala fe sobre datos ajenos | Una cuenta válida y la API |
| Persona con acceso al dispositivo | Pareja, familiar, robo | Sesión abierta en el teléfono o computador |
| Administrador o persona interna | Curiosidad, error humano | Acceso a herramientas de operación |
| Cadena de suministro | Comprometer muchas apps a la vez | Dependencias, acciones del pipeline, código sugerido por IA |

### Escenarios

| # | Escenario | Impacto | Se atiende en |
|---|---|---|---|
| A1 | Un visitante crea una cuenta con rol de administrador. | Crítico | F3 · SEC-3 — **corregido en #22** (sección 5) |
| A2 | Un usuario consulta o modifica datos de otro usuario cambiando un identificador. | Crítico | F3 · SEC-3 y F4 · SEC-4 — **corregido en `users` en #22** (sección 5); cada módulo nuevo debe repetir el control |
| A3 | Un usuario se otorga a sí mismo el rol de administrador. | Crítico | F3 · SEC-3 — **corregido en #22** (sección 5) |
| A4 | Ataque de fuerza bruta o con contraseñas filtradas contra el inicio de sesión. | Alto | F3 · SEC-3 |
| A5 | Alguien averigua si un correo está registrado. | Medio | F3 · SEC-3 |
| A6 | Se filtra la clave de firma de tokens y se fabrican sesiones. | Crítico | F2 · SEC-2 |
| A7 | Un secreto queda en el repositorio o en los logs del pipeline. | Alto | F2 · SEC-2 |
| A8 | Una dependencia vulnerable o maliciosa entra al proyecto. | Alto | F1 · SEC-1 y F2 · SEC-2 |
| A9 | Código generado por IA introduce una falla que nadie revisa. | Alto | F1 · SEC-1 |
| A10 | Se comparten datos reales de usuarios con una herramienta de IA. | Alto | F1 · SEC-1 |
| A11 | Un saldo queda mal por errores de redondeo o de conversión de moneda. | Alto | F4 · SEC-4 |
| A12 | Un movimiento se modifica sin dejar rastro. | Medio | F4 · SEC-4 |
| A13 | Una sesión robada o abierta en un dispositivo ajeno sigue activa indefinidamente. | Alto | F3 · SEC-3 y F5 · SEC-5 |
| A14 | Un error del servidor muestra información interna al usuario. | Medio | Ya mitigado en parte; se verifica en SEC-5 |
| A15 | Se pierde la base de datos sin un respaldo restaurable. | Crítico | F6 · SEC-6 |
| A16 | Abuso de la API o volumen que tumba el servicio. | Medio | F4 · SEC-4 y F6 · SEC-6 |
| A17 | Un archivo importado malicioso. | Alto | F8 · SEC-8 |

## 5. Hallazgos en el código actual

El 2026-11-14 se encontraron tres fallas de autorización en el módulo `users`. Se registraron en un GitHub Security Advisory privado, se **corrigieron el mismo día en el PR #22** y el advisory se publicó el 2026-11-14: [GHSA-j36w-x5rp-5f3w](https://github.com/davidroldan211/Finanzas-Go/security/advisories/GHSA-j36w-x5rp-5f3w).

| # | Hallazgo | Escenario | Corrección (#22) |
|---|---|---|---|
| H1 | El registro público (`POST /api/v1/users`) aceptaba el campo `role`, así que cualquiera podía crearse una cuenta de administrador. | A1 | El registro siempre crea rol `user`; el comando de creación ya no tiene rol. |
| H2 | `GET` y `PUT /api/v1/users/:id` no comprobaban que el `id` fuera el del usuario autenticado. | A2 | Solo el propio usuario o un admin; un ID ajeno responde 404, igual que uno inexistente. |
| H3 | `PUT /api/v1/users/:id` dejaba a cualquier usuario cambiar `role` (y `is_active`). | A3 | Solo un admin cambia rol o estado (403 en otro caso); el rol se valida en el dominio. |

**Verificación:** 12 pruebas automáticas nuevas y prueba de punta a punta contra la API con Postgres local (registro con `role: admin` crea `user`; leer o modificar a otro usuario da 404; ascenderse da 403; editar el propio perfil da 200).

**Lección para los módulos siguientes:** la autorización vive en la capa `application` mediante `in.Actor` (ID y rol del token). Todo caso de uso que reciba un identificador de recurso debe verificar que pertenece al actor; es el control central de SEC-4.

## 6. Criterio de aprobación de SEC-0

- [ ] Clasificación de datos revisada y aceptada.
- [ ] Lista de datos que nunca se piden aceptada.
- [ ] Normativa a evaluar asignada a alguien con fecha antes de F6.
- [ ] Escenarios A1–A17 aceptados y cada uno con la fase que lo atiende.
- [x] Hallazgos H1–H3 registrados (advisory privado) y corregidos en #22.
