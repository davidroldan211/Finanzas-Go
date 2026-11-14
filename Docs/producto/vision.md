# Visión de producto — Finanzas-Go

> Fase F0 del [roadmap](../roadmap.html). Estado: **borrador para revisión** · 2026-11-14

## Declaración de visión

Finanzas-Go ayuda a cualquier persona en Colombia a saber dónde está su dinero, en qué se le va y cuánto tiene realmente, para que pueda decidir con información en lugar de intuición.

## Problema

Una persona con ingresos propios suele tener el dinero repartido en varios lugares: efectivo, cuenta de ahorros, billetera digital, una o dos tarjetas de crédito y, cada vez más, algo en dólares. Cada banco muestra solo su parte. El resultado:

- No sabe con precisión en qué se le fue la plata del mes.
- No sabe cuánto tiene en total una vez descontadas las deudas.
- Los presupuestos y metas de ahorro viven en su cabeza o en una hoja de cálculo que abandona a las dos semanas.

## Para quién

**Usuario principal:** persona adulta en Colombia que maneja sus finanzas de forma individual.

- Tiene ingresos propios (salario, honorarios o negocio propio).
- Usa entre 2 y 6 "lugares" para su dinero: efectivo, bancos, billeteras digitales, tarjetas de crédito.
- Puede tener dinero o gastos en otra moneda (USD es el caso típico).
- No es experta en finanzas ni en contabilidad. Quiere claridad, no fórmulas.

### Tipos de usuario

| Tipo | Qué hace | Qué no puede hacer |
|---|---|---|
| **Visitante** | Ve información pública, se registra. | Nada que requiera cuenta. |
| **Usuario** | Gestiona *sus* cuentas, movimientos, presupuestos y metas. Administra su perfil y puede eliminar su cuenta. | Ver o tocar datos de otro usuario. Cambiar su propio rol. |
| **Administrador** | Opera la plataforma: soporte, desactivar cuentas abusivas, ver métricas agregadas. | **Ver los datos financieros de los usuarios.** La administración es de la plataforma, no del dinero ajeno. |

Fuera de alcance por ahora: miembros de un hogar o pareja que compartan finanzas (ver F9 del roadmap).

## Propuesta de valor

En orden de construcción:

1. **Saber en qué se va el dinero.** Registrar ingresos y gastos por categoría y verlos agrupados por mes.
2. **Ver el patrimonio total.** Todas las cuentas, tarjetas y deudas en un solo número, en la moneda principal del usuario, aunque haya cuentas en otras monedas.
3. **Controlar el presupuesto.** Fijar un límite por categoría y mes, y enterarse antes de pasarse.
4. **Ahorrar para metas.** Metas con monto, fecha y avance visible.

Las cuatro forman parte de la visión. El [alcance](alcance.md) define cuáles entran al MVP y cuáles llegan en entregas posteriores.

## Lo que Finanzas-Go no es

- **No es un banco ni una billetera.** No mueve, guarda ni transfiere dinero real. Solo registra.
- **No da asesoría financiera ni de inversión.** Muestra datos del propio usuario, no recomienda productos.
- **No se conecta a bancos en el MVP.** Los movimientos se registran a mano. La importación y la conexión llegan después.
- **No es contabilidad de empresa ni herramienta tributaria.** No genera declaraciones ni facturas.
- **No es multiusuario compartido.** Cada persona ve solo lo suyo.

## Principios de producto

1. **Registrar es rápido.** Un gasto se anota en segundos; si cuesta, el usuario deja de hacerlo y la app pierde todo su valor.
2. **Los números siempre cuadran.** Un saldo incorrecto destruye la confianza. La exactitud está por encima de las funcionalidades.
3. **Privacidad por defecto.** Se pide el mínimo de datos, nadie más ve los del usuario (tampoco los administradores) y el usuario puede llevárselos o borrarlos.
4. **Construir por partes.** Cada entrega es usable por sí sola y resuelve un problema completo antes de empezar el siguiente.
5. **Hablar como la gente.** "Tarjeta de crédito", "lo que debo", "lo que tengo"; no "pasivo corriente".

## Cómo sabremos que funciona

Indicadores iniciales, a ajustar después de la beta (F6):

| Indicador | Meta inicial |
|---|---|
| Usuarios nuevos que registran su primer movimiento el mismo día | ≥ 70 % |
| Usuarios que registran movimientos al menos 3 días por semana durante sus primeras 4 semanas | ≥ 40 % |
| Usuarios activos que vuelven a consultar su resumen mensual | ≥ 60 % |
| Reportes de saldos incorrectos | 0 sin resolver |

## Supuestos a validar

- Las personas están dispuestas a registrar movimientos a mano si hacerlo es rápido.
- Ver el patrimonio consolidado (incluyendo deudas) es suficientemente valioso para volver a la app.
- Las tasas de cambio aproximadas bastan para el patrimonio consolidado; nadie necesita precisión contable.
- El usuario colombiano típico usa billeteras digitales (Nequi, Daviplata) además de bancos.

## Decisiones tomadas

| Fecha | Decisión |
|---|---|
| 2026-11-14 | Usuario objetivo: personas en general, finanzas individuales. |
| 2026-11-14 | Mercado inicial: Colombia. Moneda principal COP, con soporte para cuentas en otras monedas desde el MVP. |
| 2026-11-14 | Entrada de movimientos en el MVP: manual. |
| 2026-11-14 | Los cuatro problemas forman parte de la visión; se construyen por partes (ver [alcance](alcance.md)). |
