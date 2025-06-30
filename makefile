# Archivo Makefile (en raíz del proyecto)

# Cargar variables desde archivo .env si existe
ifneq (,$(wildcard .env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

.PHONY: run build clean test fmt lint help

APP_NAME?=mi-proyecto
MAIN?=./cmd/finanzas


## run: Ejecuta la aplicación
run:
	go run $(APP_MAIN)

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

## lint: Linter con go vet
lint:
	go vet ./...

## fmt: Formatea el código
fmt:
	go fmt ./...

## help: Muestra los comandos disponibles
help:
	@echo "Comandos disponibles:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/^## //'

## coverage: Ejecuta tests y valida cobertura mínima
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@threshold=80.0 ; \
	actual=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//') ; \
	compare_result=$$(echo "$$actual >= $$threshold" | bc -l) ; \
	if [ "$$compare_result" -eq 1 ]; then \
		echo "✅ Cobertura suficiente ($$actual% >= $$threshold%)"; \
	else \
		echo "❌ Cobertura insuficiente ($$actual% < $$threshold%)"; \
		exit 1; \
	fi
