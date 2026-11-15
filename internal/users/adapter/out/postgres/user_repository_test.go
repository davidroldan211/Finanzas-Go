package postgres

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsUniqueViolation(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"23505 directo", &pgconn.PgError{Code: "23505"}, true},
		{"23505 envuelto", fmt.Errorf("gorm: %w", &pgconn.PgError{Code: "23505"}), true},
		{"otro código (not null)", &pgconn.PgError{Code: "23502"}, false},
		{"error que no es de Postgres", errors.New("duplicate key"), false},
		{"nil", nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isUniqueViolation(tc.err); got != tc.want {
				t.Fatalf("isUniqueViolation(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
