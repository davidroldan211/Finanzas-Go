package db

import "embed"

// Migrations embebe los archivos .sql de migraciones en el binario para que
// cmd/migrate no dependa del disco en tiempo de ejecución (funciona igual en
// desarrollo local y en CI).
//
//go:embed migrations/*.sql
var Migrations embed.FS
