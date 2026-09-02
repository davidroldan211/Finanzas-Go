package domain

import "time"

// Claims son los datos que viajan dentro de un token emitido por auth.
// Subject es un identificador de usuario opaco (string): auth no conoce
// el tipo uuid.UUID del módulo users, solo un identificador de texto.
type Claims struct {
	Subject   string
	Role      string
	ExpiresAt time.Time
}

// Token es el resultado de emitir unas Claims.
type Token struct {
	Value     string
	ExpiresAt time.Time
}
