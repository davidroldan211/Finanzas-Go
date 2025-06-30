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
	go build -o bin/$(APP_NAME) $(APP_MAIN)

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
