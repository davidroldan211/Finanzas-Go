// Comando migrate: aplica/revierte las migraciones embebidas en db.Migrations
// usando goose. cmd/ es la única capa autorizada a importar config (ver
// CLAUDE.md), igual que cmd/finanzas.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"finanzas-api/config"
	"finanzas-api/db"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	command := os.Args[1]

	// "create" escribe en el disco real (no en el FS embebido): no necesita
	// abrir conexión a la base de datos.
	if command == "create" {
		if err := runCreate(os.Args[2:]); err != nil {
			log.Fatalf("migrate: %v", err)
		}
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("migrate: cargando configuración: %v", err)
	}

	sqlDB, err := sql.Open("pgx", cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("migrate: abriendo conexión: %v", err)
	}
	defer sqlDB.Close()

	goose.SetBaseFS(db.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("migrate: fijando dialecto: %v", err)
	}

	switch command {
	case "up":
		err = goose.Up(sqlDB, migrationsDir)
	case "down":
		err = goose.Down(sqlDB, migrationsDir)
	case "status":
		err = goose.Status(sqlDB, migrationsDir)
	case "reset":
		err = goose.Reset(sqlDB, migrationsDir)
	default:
		usage()
		os.Exit(1)
	}
	if err != nil {
		log.Fatalf("migrate: %s: %v", command, err)
	}
}

func runCreate(args []string) error {
	if len(args) < 1 || args[0] == "" {
		return fmt.Errorf("uso: migrate create <nombre>")
	}
	name := args[0]

	// Sin SetBaseFS aquí: Create escribe un archivo nuevo en el disco, no en
	// el FS embebido de solo lectura. SetSequential mantiene la numeración
	// 00001, 00002... en vez del timestamp que goose usa por defecto.
	goose.SetSequential(true)
	return goose.Create(nil, "db/migrations", name, "sql")
}

func usage() {
	fmt.Println("uso: migrate <up|down|status|reset|create <nombre>>")
}
