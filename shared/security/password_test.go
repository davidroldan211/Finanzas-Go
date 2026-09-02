package security

import "testing"

func TestHashPassword_And_CheckPasswordHash_RoundTrip(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if hash == "correct-horse-battery-staple" {
		t.Fatal("HashPassword returned the plain text password")
	}

	if !CheckPasswordHash("correct-horse-battery-staple", hash) {
		t.Fatal("CheckPasswordHash rejected the correct password")
	}
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if CheckPasswordHash("wrong-password", hash) {
		t.Fatal("CheckPasswordHash accepted an incorrect password")
	}
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	if CheckPasswordHash("any-password", "not-a-valid-bcrypt-hash") {
		t.Fatal("CheckPasswordHash accepted a malformed hash")
	}
}
