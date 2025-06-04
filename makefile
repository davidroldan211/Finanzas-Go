# Archivo Makefile (en raíz del proyecto)

APP_NAME=mi-proyecto
MAIN=./cmd/finanzas

.PHONY: run build clean test fmt lint help

## Ejecuta la aplicación
run:
	go run $(MAIN)

## Compila el binario
build:
	go build -o bin/$(APP_NAME) $(MAIN)

## Ejecuta tests
test:
	go test ./...

## Limpia binarios y cache
clean:
	go clean
	rm -rf bin/

## Linter con go vet
lint:
	go vet ./...

## Formatea el código
fmt:
	go fmt ./...

## Muestra los comandos disponibles
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[36m%-12s\033[0m %s\n", $$1, $$2}'
