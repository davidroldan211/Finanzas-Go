package domain

import (
	_ "embed"
	"strings"
	"unicode/utf8"
)

const (
	// PasswordMinChars es el mínimo de caracteres (no bytes) de una
	// contraseña. OWASP ASVS 2.1.1.
	PasswordMinChars = 12
	// PasswordMaxBytes es el máximo en bytes: bcrypt ignora en silencio todo
	// lo que pase de 72 bytes, así que una contraseña más larga se rechaza.
	PasswordMaxBytes = 72
	// passwordContextMinChars es el largo mínimo que debe tener la parte
	// local del correo, el nombre o el apellido para que la contraseña no
	// pueda contenerlo. Evita rechazar por nombres muy cortos ("Al").
	passwordContextMinChars = 3
)

// Mensajes del campo "password". Los nombres evitan la palabra "password"
// a propósito: gosec (G101) toma una constante de texto con ese nombre por
// una credencial escrita en el código, y estos son solo mensajes.
const (
	msgTooShort         = "la contraseña debe tener al menos 12 caracteres"
	msgTooLong          = "la contraseña es demasiado larga"
	msgTooCommon        = "esa contraseña es muy común"
	msgContainsUserData = "la contraseña no puede contener tu correo ni tu nombre"
)

//go:embed common_passwords.txt
var commonPasswordsFile string

// commonPasswords se arma una sola vez al cargar el paquete.
var commonPasswords = parseCommonPasswords(commonPasswordsFile)

func parseCommonPasswords(file string) map[string]struct{} {
	set := make(map[string]struct{}, 10_100)
	for _, line := range strings.Split(file, "\n") {
		line = strings.TrimSuffix(line, "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		set[strings.ToLower(line)] = struct{}{}
	}
	return set
}

// ValidatePassword devuelve el mensaje del campo "password", o "" si la
// contraseña es válida. Informa solo la primera regla que falla, en este
// orden: largo mínimo, largo máximo, lista común, contexto del usuario.
// email, firstName y lastName se comparan sin distinguir mayúsculas.
func ValidatePassword(plain, email, firstName, lastName string) string {
	switch {
	case utf8.RuneCountInString(plain) < PasswordMinChars:
		return msgTooShort
	case len(plain) > PasswordMaxBytes:
		return msgTooLong
	}

	lower := strings.ToLower(plain)
	if _, common := commonPasswords[lower]; common {
		return msgTooCommon
	}

	for _, value := range passwordContext(email, firstName, lastName) {
		if strings.Contains(lower, value) {
			return msgContainsUserData
		}
	}
	return ""
}

// passwordContext devuelve, en minúsculas, los datos del usuario que la
// contraseña no puede contener: parte local del correo, nombre y apellido,
// solo los que tienen passwordContextMinChars caracteres o más.
func passwordContext(email, firstName, lastName string) []string {
	local := normalizeEmail(email)
	if at := strings.LastIndex(local, "@"); at >= 0 {
		local = local[:at]
	}

	var values []string
	for _, v := range []string{local, firstName, lastName} {
		v = strings.ToLower(strings.TrimSpace(v))
		if utf8.RuneCountInString(v) >= passwordContextMinChars {
			values = append(values, v)
		}
	}
	return values
}
