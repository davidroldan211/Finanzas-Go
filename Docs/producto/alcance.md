# Alcance y backlog — Finanzas-Go

> Fase F0 del [roadmap](../roadmap.html). Estado: **aprobado** · 2026-11-14
> Términos según el [glosario](glosario.md). Visión en [vision.md](vision.md).

## 1. Capacidades candidatas

Lista de todo lo que la app podría hacer, agrupado por épica. El ID se mantiene aunque cambie la prioridad, para que specs, ramas y PRs puedan citarlo.

| ID | Épica | Capacidad |
|---|---|---|
| CAP-01 | Identidad | Registro con verificación de correo |
| CAP-02 | Identidad | Inicio y cierre de sesión |
| CAP-03 | Identidad | Recuperar y cambiar contraseña |
| CAP-04 | Identidad | Perfil: nombre, moneda principal, zona horaria |
| CAP-05 | Identidad | Eliminar la propia cuenta y todos sus datos |
| CAP-06 | Identidad | Exportar todos los datos propios |
| CAP-10 | Cuentas | Crear, editar y archivar cuentas (efectivo, ahorros, corriente, billetera digital, tarjeta de crédito) |
| CAP-11 | Cuentas | Cada cuenta en su propia moneda (COP, USD, …) |
| CAP-12 | Cuentas | Saldo inicial al crear la cuenta |
| CAP-13 | Cuentas | Préstamos y otras deudas como cuentas |
| CAP-14 | Cuentas | Tarjeta de crédito: cupo total y cupo disponible |
| CAP-15 | Cuentas | Tarjeta de crédito: fecha de corte, fecha límite de pago y saldo a pagar del periodo |
| CAP-20 | Movimientos | Registrar ingreso o gasto |
| CAP-21 | Movimientos | Transferencia entre cuentas propias, incluso entre monedas distintas |
| CAP-22 | Movimientos | Pago de tarjeta de crédito como transferencia |
| CAP-23 | Movimientos | Corregir y anular movimientos con historial |
| CAP-24 | Movimientos | Buscar y filtrar movimientos |
| CAP-25 | Movimientos | Movimientos recurrentes |
| CAP-26 | Movimientos | Importar extractos bancarios |
| CAP-27 | Movimientos | Reglas de categorización automática |
| CAP-28 | Movimientos | Etiquetas libres y notas |
| CAP-30 | Categorías | Set predeterminado de categorías |
| CAP-31 | Categorías | Crear, renombrar y archivar categorías propias |
| CAP-32 | Categorías | Subcategorías |
| CAP-40 | Consulta | Saldo actual por cuenta |
| CAP-41 | Consulta | Patrimonio total en moneda principal |
| CAP-42 | Consulta | Resumen del mes: ingresos, gastos y gasto por categoría |
| CAP-43 | Consulta | Comparación con meses anteriores y tendencias |
| CAP-44 | Consulta | Ver y ajustar la tasa de cambio usada para consolidar |
| CAP-45 | Consulta | Tasa de cambio diaria automática desde una fuente pública |
| CAP-50 | Presupuesto | Límite mensual por categoría |
| CAP-51 | Presupuesto | Avance del presupuesto durante el mes |
| CAP-52 | Presupuesto | Aviso al acercarse o pasar el límite |
| CAP-60 | Metas | Crear meta con monto objetivo y fecha |
| CAP-61 | Metas | Registrar aportes y ver avance |
| CAP-70 | Avisos | Recordatorios de pagos próximos |
| CAP-80 | Evolución | Finanzas compartidas de hogar |
| CAP-81 | Evolución | Conexión directa con bancos |
| CAP-82 | Evolución | Asistente con IA |

## 2. Priorización

Método MoSCoW aplicado **al MVP** (v1.0). "Won't" significa "no en el MVP", no "nunca".

| Prioridad | Capacidades |
|---|---|
| **Must** — sin esto no hay MVP | CAP-01, 02, 03, 04, 10, 11, 12, 14, 15, 20, 21, 22, 23, 24, 30, 31, 40, 41, 42, 44, 45 |
| **Should** — muy deseable; entra si no retrasa el MVP | CAP-05, 06 |
| **Could** — si sobra tiempo | CAP-28, 43 |
| **Won't (en MVP)** | CAP-13, 25, 26, 27, 32, 50–52, 60, 61, 70, 80–82 |

### Por qué este corte

- Los cuatro problemas de la visión dependen de lo mismo: **movimientos bien registrados en cuentas con saldos correctos**. Sin esa base, presupuestos y metas mostrarían números falsos.
- El MVP resuelve dos problemas completos: *saber en qué se va el dinero* (CAP-42) y *ver el patrimonio total* (CAP-41).
- Multi-moneda entra desde el MVP (CAP-11, 21, 41, 44, 45) porque cambiar después cómo se guardan montos y saldos es costoso. **Esto adelanta parte de lo que el roadmap tenía en F9.**
- CAP-05 y CAP-06 son "Should" por producto, pero la normativa de datos puede volverlos obligatorios antes del lanzamiento (ver [SEC-0](../seguridad/sec-0-datos-y-amenazas.md)).

## 3. Entregas por partes

| Entrega | Resuelve | Capacidades | Fase del roadmap |
|---|---|---|---|
| **v1.0 MVP** | En qué se va el dinero + patrimonio total | Must + Should | F3, F4, F5, F6 |
| **v1.1** | Controlar el presupuesto | CAP-50, 51, 52 | F7 |
| **v1.2** | Ahorrar para metas | CAP-60, 61 | F7 |
| **v1.3** | Menos trabajo manual | CAP-25, 70, 13 | F7 |
| **v1.4** | Entender tendencias e importar datos | CAP-43, 26, 27, 32 | F8 |
| **v2.x** | Evolución | CAP-80, 81, 82 | F9 |

## 4. Mapa de historias del MVP

El recorrido del usuario va de izquierda a derecha. Debajo de cada actividad, las historias del MVP; la última fila muestra lo que llega después.

| | 1. Empezar | 2. Montar mis cuentas | 3. Registrar el día a día | 4. Ver cómo voy | 5. Mantener en orden |
|---|---|---|---|---|---|
| **v1.0** | HU-01 Registrarme | HU-05 Crear cuenta | HU-08 Registrar gasto | HU-12 Ver saldos | HU-16 Corregir movimiento |
| | HU-02 Verificar correo | HU-06 Saldo inicial | HU-09 Registrar ingreso | HU-13 Ver patrimonio | HU-17 Anular movimiento |
| | HU-03 Entrar y salir | HU-07 Cuenta en otra moneda | HU-10 Transferir entre cuentas | HU-14 Resumen del mes | HU-18 Gestionar categorías |
| | HU-04 Recuperar contraseña | HU-22 Cupo de la tarjeta | HU-11 Pagar tarjeta | HU-15 Buscar movimientos | HU-19 Archivar cuenta |
| | HU-20 Elegir moneda principal | HU-23 Corte y pago de la tarjeta | | HU-24 Ajustar tasa de cambio | HU-21 Exportar / borrar mis datos |
| **después** | | Deudas y préstamos | Recurrentes, importar extracto | Presupuestos, metas, tendencias | Reglas automáticas |

### Historias del MVP

Redacción corta. Cada una se convierte en spec completa con criterios de aceptación en F1.

| ID | Historia | Capacidad |
|---|---|---|
| HU-01 | Como visitante quiero registrarme con mi correo para empezar a usar la app. | CAP-01 |
| HU-02 | Como usuario nuevo quiero confirmar mi correo con un código para demostrar que es mío. | CAP-01 |
| HU-03 | Como usuario quiero entrar y salir de mi sesión para proteger mis datos en dispositivos compartidos. | CAP-02 |
| HU-04 | Como usuario quiero recuperar el acceso si olvido mi contraseña. | CAP-03 |
| HU-20 | Como usuario quiero elegir mi moneda principal para ver mi patrimonio en ella. | CAP-04 |
| HU-05 | Como usuario quiero crear una cuenta de tipo efectivo, banco, billetera o tarjeta para reflejar dónde está mi dinero. | CAP-10 |
| HU-06 | Como usuario quiero indicar el saldo con el que empieza una cuenta para que la app arranque con mis números reales. | CAP-12 |
| HU-07 | Como usuario quiero crear una cuenta en dólares para registrar el dinero que tengo en esa moneda. | CAP-11 |
| HU-08 | Como usuario quiero registrar un gasto con monto, categoría, cuenta y fecha en pocos segundos. | CAP-20 |
| HU-09 | Como usuario quiero registrar un ingreso para que mis saldos reflejen lo que recibo. | CAP-20 |
| HU-10 | Como usuario quiero mover dinero entre mis cuentas, aunque estén en monedas distintas (indicando cuánto salió y cuánto llegó), sin que cuente como gasto. | CAP-21 |
| HU-11 | Como usuario quiero registrar el pago de mi tarjeta de crédito para que baje mi deuda sin contar doble el gasto. | CAP-22 |
| HU-12 | Como usuario quiero ver el saldo actual de cada cuenta. | CAP-40 |
| HU-13 | Como usuario quiero ver cuánto tengo en total menos lo que debo, en mi moneda principal. | CAP-41 |
| HU-14 | Como usuario quiero ver cuánto entró, cuánto salió y en qué categorías gasté este mes. | CAP-42 |
| HU-15 | Como usuario quiero buscar y filtrar movimientos por fecha, cuenta, categoría o monto. | CAP-24 |
| HU-16 | Como usuario quiero corregir un movimiento mal registrado. | CAP-23 |
| HU-17 | Como usuario quiero anular un movimiento sin perder el rastro de que existió. | CAP-23 |
| HU-18 | Como usuario quiero crear y renombrar mis propias categorías además de las predeterminadas. | CAP-30, 31 |
| HU-19 | Como usuario quiero archivar una cuenta que ya no uso sin perder su historial. | CAP-10 |
| HU-21 | Como usuario quiero descargar o borrar todos mis datos. | CAP-05, 06 |
| HU-22 | Como usuario quiero registrar el cupo de mi tarjeta de crédito para saber cuánto cupo me queda disponible. | CAP-14 |
| HU-23 | Como usuario quiero registrar la fecha de corte y la fecha límite de pago de mi tarjeta para saber cuánto debo pagar y antes de cuándo. | CAP-15 |
| HU-24 | Como usuario quiero ver y corregir la tasa de cambio con la que se consolida mi patrimonio. | CAP-44, 45 |

## 5. Decisiones de alcance

Respuestas a las preguntas abiertas de F0 (2026-11-14). Son reglas de producto: las specs de F4 parten de aquí.

| # | Tema | Decisión | Consecuencia |
|---|---|---|---|
| 1 | Tasas de cambio | La app toma una **tasa diaria automática** de una fuente pública y el usuario puede **ajustarla**. | CAP-44 y CAP-45 pasan a Must. La fuente es un proveedor externo: entra al modelo de amenazas (A18). |
| 2 | Transferencias entre monedas | Se registra el **monto que sale y el monto que llega**; la tasa queda implícita. | HU-10 pide dos montos cuando las monedas de las cuentas difieren. |
| 3 | Tarjeta de crédito | El MVP maneja **saldo adeudado, cupo, fecha de corte y fecha límite de pago**. | Nuevas CAP-14 y CAP-15, historias HU-22 y HU-23. |
| 4 | Categorías predeterminadas | Set propuesto con categorías propias de Colombia (lista abajo). | HU-18 parte de este set. |
| 5 | Movimientos con fecha futura | **No se permiten en el MVP**: la fecha de un movimiento es hoy o anterior. | Los pagos futuros llegan con movimientos recurrentes (v1.3). |
| 6 | Edición de movimientos | Cualquier movimiento se puede **corregir o anular siempre**, con historial de cada cambio. | Los reportes de meses pasados pueden cambiar; el historial explica por qué. |

### Categorías predeterminadas

| Tipo | Categorías |
|---|---|
| Gasto | Vivienda y arriendo, Mercado, Transporte, Servicios públicos, Salud, Educación, Restaurantes, Entretenimiento, Compras, Deudas e intereses, Otros gastos |
| Ingreso | Salario, Honorarios, Prima, Cesantías, Otros ingresos |
