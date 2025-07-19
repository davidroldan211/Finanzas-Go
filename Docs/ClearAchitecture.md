# Clean Architecture 

La **Clean Architecture**, propuesta por Robert C. Martin (Uncle Bob), busca separar los componentes de un sistema en capas independientes, de forma que cada una tenga una responsabilidad clara y el flujo de dependencias siempre vaya desde el exterior hacia el centro.

![Texto alternativo](a1f9e6c33c27b9695b45d98ae6f7708b.jpg)

---

## DOMAIN
**Contiene las reglas de negocio puras** del sistema. Es el núcleo de la arquitectura.

Elementos clave:
- **Entities**: Objetos del dominio con identidad propia y lógica de negocio.
- **Value Objects**: Objetos sin identidad, inmutables, que representan conceptos.
- **Domain Events**: Eventos que representan algo que ocurrió dentro del dominio.
- **Enums & Exceptions**: Tipos fuertes y errores específicos del dominio.
- **Repositories (Interfaces)**: Contratos para acceso a datos.

No depende de ninguna tecnología o librería externa.

---

## APPLICATION
Encapsula los **casos de uso del sistema**. Orquesta las entidades del dominio para cumplir una acción.

 Elementos clave:
- **Use Cases**: Lógica específica de aplicación (ej: registrar usuario).
- **Application Services**: Coordinan el uso de entidades/repositorios.
- **Commands / Queries**: Modelos para entrada/salida (estilo CQRS).
- **External Interfaces**: Interfaces hacia infraestructura (ej: repositorios, servicios).

Aquí se toman decisiones de flujo y validación de negocio, sin saber nada de frameworks.

---

## INFRASTRUCTURE
Se encarga de **las dependencias externas** al sistema.

Elementos clave:
- **Databases / Repository Implementations**: Implementación de acceso a datos.
- **HTTP Clients / Email / Cloud Storage**: Servicios externos.
- **Message Brokers**: Comunicación asíncrona.
- **Identity Providers**: Autenticación/autorización.

Esta capa puede cambiar sin tocar el núcleo del sistema gracias a las interfaces.

---

## PRESENTATION
Es el **punto de entrada** al sistema. Define cómo se reciben y responden las solicitudes.

Elementos clave:
- **API Endpoints / GraphQL / gRPC**: Interfaces expuestas al cliente.
- **ASP.NET / Middleware**: Frameworks y componentes del entorno.
- **Services (DI)**: Configuración de inyección de dependencias.
- **Exceptions**: Manejo de errores presentables.

Actúa como "composición raíz" del sistema, donde se conectan los casos de uso con los adaptadores reales.



