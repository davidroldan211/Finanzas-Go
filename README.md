# 🏦 Finanzas Personales

Sistema backend para la gestión de finanzas personales, desarrollado en Go con una arquitectura hexagonal, modular y escalable, orientado a entornos empresariales.

---

## 📁 Estructura del Proyecto

El proyecto esta construido siguiendo los lineamientos de [Arquitectura Hexagonal (Puertos y Adaptadores)](docs/HexagonalArchitecture.md)

## ⚙️ Requisitos

- Go 1.21 o superior
- Docker y Docker Compose
- `make` 


## 🔧 Instalación

```bash
git clone https://github.com/davidroldan211/Finanzas-Go.git
```
```bash
cd finanzas-personales
```
```bash
go mod tidy
```

## 🐳 Docker
```bash
docker compose up -d
```
Este comando levantará la infraestructura definida (como PostgreSQL, Redis, etc.).

## 🗃️ Migraciones de base de datos

Las tablas **no** se crean solas (no hay `AutoMigrate`): hay que aplicar las migraciones antes de correr la app por primera vez.

```bash
make migrate-up
```

### Cómo funciona

- Los archivos `.sql` viven en `db/migrations/`, numerados en orden (`00001_...`, `00002_...`, ...) y se aplican en ese orden.
- Cada archivo tiene un bloque `-- +goose Up` (qué crea) y uno `-- +goose Down` (cómo deshacerlo).
- `cmd/migrate` (usa [goose](https://github.com/pressly/goose)) los aplica contra la base configurada en `.env`. Los `.sql` quedan embebidos en el binario (`db/migrations.go`), así que no dependen de encontrar el directorio en disco.
- Postgres registra qué migraciones ya corrieron, así que `make migrate-up` es seguro de repetir: solo aplica las pendientes.

### Comandos

```bash
make migrate-up                        # aplica las migraciones pendientes
make migrate-status                    # muestra cuáles están aplicadas y cuáles no
make migrate-down                      # revierte la última migración aplicada
make migrate-reset                     # revierte todas (borra las tablas)
make migrate-create name=nombre_algo   # crea db/migrations/000NN_nombre_algo.sql
```

### Regla al agregar un cambio de esquema

1. `make migrate-create name=lo_que_cambia` — crea el archivo con el siguiente número.
2. Completar `-- +goose Up` con el cambio y `-- +goose Down` con cómo revertirlo.
3. **Nunca editar una migración que ya se mergeó a `main`.** Si algo quedó mal, se corrige con una migración nueva (número mayor), no editando la vieja — de lo contrario una base ya creada con el archivo original queda con un esquema que el `.sql` ya no describe, y no hay forma de notarlo.

## ▶️ Ejecución
Modo local
```bash
go run ./cmd/finanzas
```
Usando Makefile
```bash
make run       # Ejecutar app
```
```bash
make build     # Compilar binario
```
```bash
make test      # Ejecutar tests
```


## 🧪 Pruebas
```bash
go test ./...
```


## 🧰 Comandos Útiles
```bash
make run       # Ejecutar el servicio
```
```bash
make build     # Compilar binario
```
```bash
make clean     # Limpiar binarios
```
```bash
make test      # Ejecutar tests
```
```bash
make fmt       # Formatear código
```
```bash
make lint      # Revisar calidad del código con go vet
```
```bash 
make coverage  # Valida cobertura de pruebas 
```
```bash
make migrate-up      # Aplicar migraciones pendientes
```
```bash
make migrate-status  # Ver estado de las migraciones
```
## 📦 Tecnologías Utilizadas
- Go (Golang)
- Arquitectura Hexagonal (Puertos y Adaptadores)
- Docker / Docker Compose
- PostgreSQL
- (Opcional) Gin / Echo / Fiber como framework web
- GORM (ORM para Go)


## 🛠 Buenas Prácticas Aplicadas
- Estructura basada en cmd/ para binarios
- internal/ para encapsular lógica de negocio
- Separación por puertos y adaptadores: domain, port, application, adapter
- Automatización con Makefile
- Uso de variables de entorno en config/



## 🧩 Roadmap
- Conexión a base de datos
- Autenticación con JWT
- CRUD de presupuestos personales
- Reportes mensuales
- Dashboard financiero



## 🤝 Contribuciones
1. Haz un fork del repositorio.

2. Crea una rama con tu funcionalidad: git checkout -b feature/nueva-funcionalidad.

3. Realiza tus cambios y haz commit: git commit -m 'feat: agrega nueva funcionalidad'.

4. Haz push a la rama: git push origin feature/nueva-funcionalidad.

5. Abre un Pull Request.



## 📝 Licencia
MIT © 2025 David Roldán

