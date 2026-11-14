# Glosario del dominio — Finanzas-Go

> Fase F0 del [roadmap](../roadmap.html). Estado: **aprobado** · 2026-11-14

Un solo significado por término. Specs, código, interfaz y conversaciones con la IA usan estas palabras con este sentido. Si un término nuevo aparece en una spec, se agrega aquí primero.

La columna **En la interfaz** indica cómo se le muestra al usuario cuando el término técnico no es natural.

## Personas

| Término | Definición | En la interfaz |
|---|---|---|
| **Usuario** | Persona registrada que gestiona sus propias finanzas. Es dueña de todos sus datos. | — |
| **Administrador** | Persona que opera la plataforma. Puede gestionar usuarios (por ejemplo, desactivarlos), pero no ve sus datos financieros. | — |
| **Visitante** | Persona que aún no tiene cuenta o no ha iniciado sesión. | — |

## Dinero y monedas

| Término | Definición | En la interfaz |
|---|---|---|
| **Moneda** | Unidad monetaria identificada por su código ISO 4217 (COP, USD, EUR). | — |
| **Moneda principal** | Moneda que el usuario elige para ver totales consolidados. Por defecto, COP. | "Mi moneda" |
| **Monto** | Cantidad de dinero en una moneda específica. Un monto sin moneda no existe. | — |
| **Tasa de cambio** | Cuántas unidades de la moneda principal equivalen a una unidad de otra moneda en una fecha dada. La app toma una tasa diaria de una fuente pública y el usuario puede ajustarla. Solo se usa para consolidar; nunca cambia los montos registrados. | — |
| **Tasa ajustada** | Tasa de cambio que el usuario corrigió a mano para una fecha; reemplaza a la automática en sus totales. | — |

## Cuentas

| Término | Definición | En la interfaz |
|---|---|---|
| **Cuenta** | Lugar donde el usuario tiene dinero o una deuda. Tiene nombre, tipo, moneda y saldo. **No** es la cuenta de usuario para iniciar sesión. | "Mis cuentas" |
| **Tipo de cuenta** | Clasificación de una cuenta: efectivo, cuenta de ahorros, cuenta corriente, billetera digital, tarjeta de crédito. (Préstamo llega después del MVP.) | — |
| **Cuenta de activo** | Cuenta cuyo saldo es dinero que el usuario tiene: efectivo, ahorros, corriente, billetera digital. | "Lo que tengo" |
| **Cuenta de pasivo** | Cuenta cuyo saldo es dinero que el usuario debe: tarjeta de crédito, préstamo. | "Lo que debo" |
| **Saldo inicial** | Saldo de la cuenta en el momento en que el usuario empieza a registrarla. Punto de partida del saldo. | — |
| **Saldo** | Saldo inicial más todos los movimientos vigentes de la cuenta. Siempre en la moneda de la cuenta. | — |
| **Cupo** | Monto máximo que el banco permite deber en una tarjeta de crédito. | — |
| **Cupo disponible** | Cupo menos el saldo adeudado de la tarjeta. | — |
| **Fecha de corte** | Día del mes en que el banco cierra el periodo de la tarjeta y calcula lo que hay que pagar. | — |
| **Fecha límite de pago** | Último día para pagar lo facturado al corte sin cargos de mora. | — |
| **Saldo a pagar del periodo** | Lo que se gastó con la tarjeta entre el corte anterior y el último corte, menos lo que ya se pagó. | "Lo que debo pagar" |
| **Cuenta archivada** | Cuenta que ya no se usa. Conserva su historial y sigue contando en reportes pasados, pero no aparece para registrar movimientos nuevos. | — |
| **Patrimonio** | Suma de los saldos de las cuentas de activo menos la suma de los saldos de las cuentas de pasivo, convertida a la moneda principal. | "Mi patrimonio" o "Total" |

## Movimientos

| Término | Definición | En la interfaz |
|---|---|---|
| **Movimiento** | Registro de dinero que entra, sale o se mueve. Tiene monto, fecha (hoy o anterior; en el MVP no hay fechas futuras), cuenta, descripción opcional y, si aplica, categoría. Es el dato central de la app. | — |
| **Ingreso** | Movimiento que aumenta el dinero del usuario desde fuera (salario, pago recibido). Lleva categoría de ingreso. | — |
| **Gasto** | Movimiento que disminuye el dinero del usuario hacia fuera (compra, servicio). Lleva categoría de gasto. Una compra con tarjeta de crédito es un gasto en la cuenta de la tarjeta. | — |
| **Transferencia** | Movimiento de dinero entre dos cuentas del mismo usuario. No es ingreso ni gasto y no lleva categoría. Si las cuentas tienen monedas distintas, registra el monto que sale y el que llega; la tasa queda implícita. | — |
| **Pago de tarjeta** | Transferencia desde una cuenta de activo hacia una tarjeta de crédito. Reduce la deuda; el gasto ya se contó cuando se usó la tarjeta. | — |
| **Corrección** | Cambio a un movimiento existente, sin importar su antigüedad. Queda registro de qué cambió y cuándo. | "Editar" |
| **Anulación** | Marcar un movimiento como no válido. Deja de afectar saldos y reportes, pero no se borra. | "Anular" |
| **Movimiento recurrente** | Plantilla que genera movimientos con una periodicidad (mensual, quincenal…). Posterior al MVP. | — |

## Categorías

| Término | Definición | En la interfaz |
|---|---|---|
| **Categoría** | Clasificación de un ingreso o un gasto según su propósito (Mercado, Transporte, Salario). Cada categoría es de ingreso o de gasto, no de ambos. | — |
| **Categoría predeterminada** | Categoría que la app crea para cada usuario nuevo. El usuario puede renombrarla o archivarla. | — |
| **Categoría propia** | Categoría creada por el usuario. | — |
| **Subcategoría** | Categoría dentro de otra. Posterior al MVP. | — |

## Consulta y planeación

| Término | Definición | En la interfaz |
|---|---|---|
| **Periodo** | Rango de fechas para agrupar movimientos. En el MVP, el mes calendario. | — |
| **Resumen mensual** | Total de ingresos, total de gastos, diferencia y gasto por categoría de un periodo, en la moneda principal. | "Mi mes" |
| **Presupuesto** | Monto máximo que el usuario planea gastar en una categoría durante un periodo. Posterior al MVP (v1.1). | — |
| **Meta de ahorro** | Monto objetivo con fecha límite opcional que el usuario quiere reunir. Posterior al MVP (v1.2). | "Meta" |
| **Aporte** | Monto que el usuario asigna a una meta de ahorro. | — |

## Términos que no usamos

| No usar | Usar | Motivo |
|---|---|---|
| Transacción | Movimiento | "Transacción" se confunde con transacción de base de datos. |
| Débito / crédito (contables) | Gasto / ingreso, o "sale" / "entra" | El usuario no es contador; y "crédito" choca con "tarjeta de crédito". |
| Balance | Saldo o patrimonio | Ambiguo. |
| Borrar movimiento | Anular movimiento | Los movimientos no se borran; se anulan. |
