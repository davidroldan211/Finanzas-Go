package postgres

import "github.com/google/uuid"

// userCredentialsModel es una proyección read-only de la tabla "users",
// cuyo esquema es fuente de verdad en db/migrations/users.sql (propiedad
// del módulo users). auth solo lee estas 5 columnas y nunca escribe en
// esta tabla: no importa el modelo de persistencia de users para no
// acoplar los dos módulos.
type userCredentialsModel struct {
	ID       uuid.UUID
	Email    string
	Password string
	Role     string
	IsActive bool
}

func (userCredentialsModel) TableName() string { return "users" }
