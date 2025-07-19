# 🏦 Finanzas Personales

Sistema backend para la gestión de finanzas personales, desarrollado en Go con una arquitectura limpia, modular y escalable, orientado a entornos empresariales.

---

## 📁 Estructura del Proyecto

El proyecto esta contruido siguiendo los lineamientos de Clean Architecture 
[Clean Architecture ](Docs/ClearAchitecture.md)

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

## 📦 Tecnologías Utilizadas
- Go (Golang)
- Clean Architecture
- Docker / Docker Compose
- PostgreSQL
- (Opcional) Gin / Echo / Fiber como framework web
- GORM (ORM para Go)


## 🛠 Buenas Prácticas Aplicadas
- Estructura basada en cmd/ para binarios
- internal/ para encapsular lógica de negocio
- Separación por capas: dominio, casos de uso, repositorio, handlers
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

