# 🏦 Finanzas Personales

[![CI](https://github.com/davidroldan211/Finanzas-Go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/davidroldan211/Finanzas-Go/actions/workflows/ci.yml)
[![Seguridad](https://github.com/davidroldan211/Finanzas-Go/actions/workflows/security.yml/badge.svg?branch=main)](https://github.com/davidroldan211/Finanzas-Go/actions/workflows/security.yml)

Sistema backend para la gestión de finanzas personales, desarrollado en Go con una arquitectura hexagonal, modular y escalable, orientado a entornos empresariales.

---

## 📁 Estructura del Proyecto

El proyecto esta construido siguiendo los lineamientos de [Arquitectura Hexagonal (Puertos y Adaptadores)](docs/HexagonalArchitecture.md)

## ⚙️ Requisitos

- Go 1.26 o superior (la versión exacta está en `go.mod`)
- Docker y Docker Compose
- `make` 


## 🔧 Instalación

```bash
git clone https://github.com/davidroldan211/Finanzas-Go.git
```
```bash
cd Finanzas-Go
```
```bash
go mod tidy
```

## 🐳 Docker
```bash
docker compose up -d
```
Este comando levanta PostgreSQL 18.

### Si ya tenías la base con PostgreSQL 14

Un volumen creado con Postgres 14 no sirve con Postgres 18: el contenedor no arranca. Con `docker compose` todavía en la versión 14 (rama o commit anterior), exporta los datos:

```bash
docker compose exec -T db pg_dumpall -U "$DB_USER" > respaldo-pg14.sql
```

Luego, ya con la versión 18, recrea el volumen y restaura. **`down -v` borra el volumen con los datos de Postgres 14**: confirma que `respaldo-pg14.sql` tiene contenido antes de correrlo.

```bash
docker compose down -v
```
```bash
docker compose up -d
```
```bash
docker compose exec -T db psql -U "$DB_USER" -d postgres < respaldo-pg14.sql
```

Si los datos locales no importan, basta con `docker compose down -v` y `make migrate-up`.

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
make lint      # go vet + fronteras de la arquitectura
```
```bash
make coverage  # Pruebas con -race y cobertura mínima
```
```bash
make ci        # Mismos controles que el pipeline
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
El plan por fases, con puntos de control de seguridad, está en [Docs/roadmap.html](Docs/roadmap.html). Visión y alcance del producto en [Docs/producto/](Docs/producto/).


## 🤝 Contribuciones
Ver [CONTRIBUTING.md](CONTRIBUTING.md). Las vulnerabilidades se reportan en privado según [SECURITY.md](SECURITY.md).



## 📝 Licencia
© 2025 David Roldán. Todos los derechos reservados. El código es visible públicamente, pero no se concede licencia de uso, copia ni distribución.

