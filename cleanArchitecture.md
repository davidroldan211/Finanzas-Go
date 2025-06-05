# Clean Architecture en Go

## ¿Qué es Clean Architecture?

Clean Architecture es un patrón de diseño de software creado por Robert C. Martin (Uncle Bob) que organiza el código en capas concéntricas, donde las dependencias siempre apuntan hacia adentro. El objetivo principal es construir aplicaciones:

- Independientes de frameworks
- Fácilmente testeables
- Independientes de la UI
- Independientes de la base de datos
- Independientes de servicios externos

---

## Diagrama de Capas

```
+-----------------------------+
|     Frameworks & Drivers   |  -> HTTP, DB, CLI, etc.
+-----------------------------+
| Interface Adapters         |  -> Controllers, Presenters, Gateways
+-----------------------------+
| Application Business Rules |  -> Use Cases
+-----------------------------+
| Enterprise Business Rules  |  -> Entities
+-----------------------------+
```

Las dependencias deben fluir desde las capas externas hacia las internas.

---

## Estructura de Carpetas

```
proyecto/
├── cmd/                    # Punto de entrada de la aplicación
├── internal/              # Lógica del dominio y de aplicación
│   ├── domain/           # Entidades y contratos (Enterprise Rules)
│   ├── usecases/         # Casos de uso (Application Rules)
│   ├── adapters/         # Adaptadores de interfaz (Controllers, Repos)
│   └── infrastructure/   # Implementaciones externas (DB, Web, etc)
├── pkg/                  # Librerías reutilizables
└── docs/                 # Documentación
```

---

## 📁 Descripción de Capas y Responsabilidades

### 1. `domain/` - Enterprise Business Rules
**Responsabilidad:** Modelar las entidades del negocio y sus reglas.

Contiene:
- Entidades puras del negocio
- Validaciones de dominio
- Contratos (interfaces) de repositorios

Ejemplo: `User`, `Product`, `Order`, etc.

#### Ejemplo:
```go
// internal/usecases/user_usecase.go
package usecases

import (
    "errors"
    "project/internal/domain"
)

type UserRepository interface {
    Save(user domain.User) error
}

type UserUseCase struct {
    repo UserRepository
}

func NewUserUseCase(r UserRepository) *UserUseCase {
    return &UserUseCase{repo: r}
}

func (uc *UserUseCase) Register(email, name string) error {
    user := domain.User{Email: email, Name: name}
    if !user.IsValid() {
        return errors.New("invalid user")
    }
    return uc.repo.Save(user)
}
```

### 2. `usecases/` - Application Business Rules
**Responsabilidad:** Coordinar la ejecución de acciones del negocio.

Contiene:
- Lógica de casos de uso
- Orquestación entre entidades y adaptadores
- Validaciones específicas de la aplicación

Ejemplo: `RegisterUser`, `DeleteOrder`, `UpdateProfile`

#### Ejemplo:
```go
// internal/usecases/user_usecase.go
package usecases

import (
    "errors"
    "project/internal/domain"
)

type UserRepository interface {
    Save(user domain.User) error
}

type UserUseCase struct {
    repo UserRepository
}

func NewUserUseCase(r UserRepository) *UserUseCase {
    return &UserUseCase{repo: r}
}

func (uc *UserUseCase) Register(email, name string) error {
    user := domain.User{Email: email, Name: name}
    if !user.IsValid() {
        return errors.New("invalid user")
    }
    return uc.repo.Save(user)
}
```

### 3. `adapters/` - Interface Adapters
**Responsabilidad:** Adaptar datos entre el mundo exterior y la lógica interna.

Contiene:
- Controladores HTTP
- Presentadores/serializadores
- Implementaciones concretas de interfaces (repositorios, gateways)

Ejemplo: `UserController`, `UserRepositoryMySQL`, `UserJSONPresenter`

#### Ejemplo:
```go
// internal/adapters/http/user_handler.go
package http

import (
    "encoding/json"
    "net/http"
    "project/internal/usecases"
)

type UserHandler struct {
    UC *usecases.UserUseCase
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
        Name  string `json:"name"`
    }
    _ = json.NewDecoder(r.Body).Decode(&req)
    err := h.UC.Register(req.Email, req.Name)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
}
```

### 4. `infrastructure/` - Frameworks & Drivers
**Responsabilidad:** Implementaciones técnicas concretas, como DB, routers, etc.

Contiene:
- Configuración de base de datos
- Inicialización de servidor HTTP
- Conexiones externas (APIs, correo, archivos, etc)

Ejemplo: `MySQLConnection`, `GinRouter`, `MailProvider`

#### Ejemplo:
```go
// internal/infrastructure/mysql/user_repository.go
package mysql

import (
    "database/sql"
    "project/internal/domain"
)

type MySQLUserRepository struct {
    DB *sql.DB
}

func (r *MySQLUserRepository) Save(user domain.User) error {
    _, err := r.DB.Exec("INSERT INTO users (email, name) VALUES (?, ?)", user.Email, user.Name)
    return err
}
```


---

## 🔄 Flujo de Dependencias

```
HTTP Request
    ↓
Controller (adapters)
    ↓
UseCase (usecases)
    ↓
Entity / Repository Interface (domain)
    ↓
Repository Implementation (adapters) → DB (infrastructure)
```

**Regla de Oro:** Ninguna capa interna debe conocer detalles de capas externas.

---

## ✅ Beneficios de Clean Architecture

- **Separación de responsabilidades** clara
- **Alto nivel de testabilidad**
- **Fácil mantenimiento y extensión**
- **Independencia tecnológica** (puedes cambiar DB o framework sin afectar el core)

---

## 📅 Ejemplo de Caso de Uso

```go
func (uc *UserUseCase) CreateUser(email, name string) (*entities.User, error) {
    existing, _ := uc.repo.FindByEmail(email)
    if existing != nil {
        return nil, errors.New("email ya registrado")
    }

    user := &entities.User{
        Email: email,
        Name: name,
        CreatedAt: time.Now(),
    }

    if !user.IsValidEmail() {
        return nil, errors.New("email no válido")
    }

    err := uc.repo.Save(user)
    return user, err
}
```

---

## 🎓 Recursos Adicionales

- [The Clean Architecture - Uncle Bob](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design)
- [Go Clean Architecture Examples](https://github.com/larrabee/go-clean-architecture)
