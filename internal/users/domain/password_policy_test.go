package domain

import (
	"strings"
	"testing"
)

const (
	testEmail     = "ana.perez@ejemplo.co"
	testFirstName = "Ana"
	testLastName  = "Pérez"
)

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name      string
		plain     string
		email     string
		firstName string
		lastName  string
		want      string
	}{
		{"11 caracteres falla", "rio-cometa-", testEmail, testFirstName, testLastName, msgTooShort},
		{"12 caracteres pasa", "rio-cometa-z", testEmail, testFirstName, testLastName, ""},
		{"73 bytes falla", strings.Repeat("x", 73), testEmail, testFirstName, testLastName, msgTooLong},
		{"72 bytes pasa", strings.Repeat("x", 72), testEmail, testFirstName, testLastName, ""},
		{"18 caracteres multibyte (36 bytes) pasa", strings.Repeat("ñ", 18), testEmail, testFirstName, testLastName, ""},
		{"11 caracteres multibyte falla aunque sean 22 bytes", strings.Repeat("ñ", 11), testEmail, testFirstName, testLastName, msgTooShort},
		{"lista común (SecLists) falla", "unbelievable", testEmail, testFirstName, testLastName, msgTooCommon},
		{"lista común en mayúsculas falla", "UNBELIEVABLE", testEmail, testFirstName, testLastName, msgTooCommon},
		{"lista común en español falla", "Contraseña123", testEmail, testFirstName, testLastName, msgTooCommon},
		{"contiene la parte local del correo falla", "x-ANA.PEREZ-cometa", testEmail, testFirstName, testLastName, msgContainsUserData},
		{"contiene el nombre falla", "cometa-ana-azul", testEmail, testFirstName, testLastName, msgContainsUserData},
		{"contiene el apellido falla", "cometa-PÉREZ-azul", testEmail, testFirstName, testLastName, msgContainsUserData},
		{"nombre de 2 caracteres no cuenta", "al-cometa-azul", "x9@ejemplo.co", "Al", "Li", ""},
		{"dominio del correo no cuenta", "ejemplo.co-cometa", testEmail, testFirstName, testLastName, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ValidatePassword(tc.plain, tc.email, tc.firstName, tc.lastName); got != tc.want {
				t.Fatalf("ValidatePassword(%q) = %q, want %q", tc.plain, got, tc.want)
			}
		})
	}
}

func TestCommonPasswords_LoadedFromEmbeddedFile(t *testing.T) {
	if n := len(commonPasswords); n < 10_000 {
		t.Fatalf("expected at least 10000 common passwords, got %d", n)
	}
	for key := range commonPasswords {
		if strings.HasPrefix(key, "#") {
			t.Fatalf("comment line loaded as a password: %q", key)
		}
	}
}

func TestCommonPasswordsFile_CitesSourceAndLicense(t *testing.T) {
	for _, want := range []string{"SecLists", "MIT"} {
		if !strings.Contains(commonPasswordsFile, want) {
			t.Fatalf("common_passwords.txt header must mention %q", want)
		}
	}
}
