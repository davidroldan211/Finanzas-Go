# Archivo Makefile (en raíz del proyecto).

# Cargar variables desde archivo .env si existe
ifneq (,$(wildcard .env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

.PHONY: run build clean test fmt fmt-check tidy-check vet lint golangci vuln ci arch help coverage migrate-up migrate-down migrate-status migrate-reset migrate-create

APP_NAME?=mi-proyecto
MAIN?=./cmd/finanzas


## run: Ejecuta la aplicación
run:
	go run $(MAIN)

## build: Compila el binario
build:
	@if [ ! -d "$(MAIN)" ]; then \
		echo "❌ ERROR: El directorio $(MAIN) no existe."; \
		exit 1; \
	fi
	@if ! ls $(MAIN)/*.go > /dev/null 2>&1; then \
		echo "❌ ERROR: No hay archivos Go en $(MAIN)."; \
		exit 1; \
	fi
	@mkdir -p bin
	go build -o bin/$(APP_NAME) $(MAIN)

## test: Ejecuta tests
test:
	go test ./...

## clean: Limpia binarios y cache
clean:
	go clean
	rm -rf bin/

## vet: Análisis estático estándar de Go
vet:
	go vet ./...

## lint: go vet + fronteras de la arquitectura
lint: vet arch

## golangci: Linter completo (requiere golangci-lint v2 instalado)
golangci:
	golangci-lint run ./...

## vuln: Busca vulnerabilidades conocidas en dependencias (requiere govulncheck)
vuln:
	govulncheck ./...

## fmt-check: Falla si hay archivos sin formatear (no modifica nada)
fmt-check:
	@out=$$(gofmt -l $$(git ls-files '*.go')); \
	if [ -n "$$out" ]; then \
		echo "❌ Archivos sin formato (correr make fmt):"; echo "$$out"; exit 1; \
	fi
	@echo "✅ Formato correcto"

## tidy-check: Falla si go.mod/go.sum no están al día (no modifica nada)
tidy-check:
	go mod tidy -diff

## ci: Ejecuta localmente los controles del pipeline (sin Postgres ni golangci)
ci: fmt-check tidy-check lint coverage build

## arch: Verifica que el núcleo hexagonal (domain/application/port) no
## dependa de frameworks ni de infraestructura
arch:
	@echo "🏛️  Verificando fronteras de la arquitectura hexagonal..."
	@! grep -RIn --include='*.go' \
		-E '"(github.com/gin-gonic|gorm.io|finanzas-api/config|finanzas-api/internal/httpx)' \
		internal/*/domain internal/*/application internal/*/port 2>/dev/null \
		| grep -v '_test.go' \
		|| (echo "\033[1;31m❌ El núcleo (domain/application/port) importa infraestructura\033[0m"; exit 1)
	@echo "\033[1;32m✅ Núcleo limpio de infraestructura\033[0m"

## fmt: Formatea el código
fmt:
	go fmt ./...

## help: Muestra los comandos disponibles
help:
	@echo "Comandos disponibles:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/^## //'

## coverage: Ejecuta tests y valida cobertura mínima
coverage:
	@echo ""
	@echo "📦 Ejecutando pruebas y generando perfil de cobertura... "
	@echo ""
	@go test -race -coverprofile=coverage.out ./...

	@echo ""
	@echo "📊 Resumen de cobertura por función:"
	@echo ""
	@go tool cover -func=coverage.out

	@echo ""
	@echo "🔍 Verificando umbral mínimo de cobertura..."
	@echo ""
	@threshold=35.0 ; \
	actual=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//') ; \
	compare_result=$$(echo "$$actual >= $$threshold" | bc -l) ; \
	if [ "$$compare_result" -eq 1 ]; then \
		echo ""; \
		echo "\033[1;32m✅ Cobertura suficiente:\033[0m $$actual% >= $$threshold%" ; \
		echo ""; \
	else \
		echo ""; \
		echo "\033[1;31m❌ Cobertura insuficiente:\033[0m Actual = $$actual% < Deseado = $$threshold%" ; \
		echo ""; \
		exit 1 ; \
	fi

## migrate-up: Aplica las migraciones pendientes
migrate-up:
	go run ./cmd/migrate up

## migrate-down: Revierte la última migración
migrate-down:
	go run ./cmd/migrate down

## migrate-status: Muestra qué migraciones están aplicadas
migrate-status:
	go run ./cmd/migrate status

## migrate-reset: Revierte todas las migraciones
migrate-reset:
	go run ./cmd/migrate reset

## migrate-create: Crea una migración nueva (make migrate-create name=add_email_verified)
migrate-create:
	go run ./cmd/migrate create $(name)

