package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewUser_Valid(t *testing.T) {
	u, err := NewUser("  A@B.COM  ", "  Ana  ", "  Pérez  ", "hashed", "")
	if err != nil {
		t.Fatalf("NewUser returned error: %v", err)
	}
	if u.Email != "a@b.com" {
		t.Fatalf("expected normalized email 'a@b.com', got %q", u.Email)
	}
	if u.FirstName != "Ana" || u.LastName != "Pérez" {
		t.Fatalf("expected trimmed names, got %q %q", u.FirstName, u.LastName)
	}
	if u.Role != RoleUser {
		t.Fatalf("expected default role %q, got %q", RoleUser, u.Role)
	}
	if !u.IsActive {
		t.Fatal("expected new user to be active by default")
	}
}

func TestNewUser_ReportsAllInvalidFields(t *testing.T) {
	_, err := NewUser("", "", "", "hashed", "")

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T: %v", err, err)
	}
	for _, field := range []string{"email", "first_name", "last_name"} {
		if _, ok := ve.Fields[field]; !ok {
			t.Errorf("expected ValidationError to report field %q, got %v", field, ve.Fields)
		}
	}
}

func TestNewUser_InvalidEmailFormat(t *testing.T) {
	_, err := NewUser("not-an-email", "Ana", "Pérez", "hashed", "")

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T: %v", err, err)
	}
	if _, ok := ve.Fields["email"]; !ok {
		t.Fatalf("expected email field error, got %v", ve.Fields)
	}
}

func TestApplyUpdate_PartialChangesOnly(t *testing.T) {
	u, err := NewUser("a@b.com", "Ana", "Pérez", "hashed", RoleUser)
	if err != nil {
		t.Fatalf("NewUser returned error: %v", err)
	}

	newFirstName := "Ana María"
	if err := u.ApplyUpdate(nil, &newFirstName, nil, nil, nil); err != nil {
		t.Fatalf("ApplyUpdate returned error: %v", err)
	}

	if u.Email != "a@b.com" {
		t.Fatalf("expected email untouched, got %q", u.Email)
	}
	if u.FirstName != "Ana María" {
		t.Fatalf("expected first name updated, got %q", u.FirstName)
	}
	if u.LastName != "Pérez" {
		t.Fatalf("expected last name untouched, got %q", u.LastName)
	}
}

func TestApplyUpdate_InvalidResultRejected(t *testing.T) {
	u, err := NewUser("a@b.com", "Ana", "Pérez", "hashed", RoleUser)
	if err != nil {
		t.Fatalf("NewUser returned error: %v", err)
	}

	empty := ""
	if err := u.ApplyUpdate(&empty, nil, nil, nil, nil); err == nil {
		t.Fatal("expected ApplyUpdate to reject an empty email")
	}
	if u.Email != "a@b.com" {
		t.Fatalf("expected user unchanged after rejected update, got email %q", u.Email)
	}
}

func TestGetFullName(t *testing.T) {
	u := &User{FirstName: "Ana", LastName: "Pérez"}
	if got := u.GetFullName(); got != "Ana Pérez" {
		t.Fatalf("expected 'Ana Pérez', got %q", got)
	}
}

func TestIsValidForAuth(t *testing.T) {
	deletedAt := time.Now()
	cases := []struct {
		name     string
		user     User
		expected bool
	}{
		{"active and not deleted", User{IsActive: true, DeletedAt: nil}, true},
		{"inactive and not deleted", User{IsActive: false, DeletedAt: nil}, false},
		{"active but deleted", User{IsActive: true, DeletedAt: &deletedAt}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.user.IsValidForAuth(); got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
