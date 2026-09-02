package security

import "testing"

func TestBcryptHasher_Hash_And_Matches_RoundTrip(t *testing.T) {
	h := NewBcryptHasher(4) // costo mínimo válido: tests rápidos

	hash, err := h.Hash("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}
	if hash == "" || hash == "correct-horse-battery-staple" {
		t.Fatalf("Hash returned an invalid hash: %q", hash)
	}

	if !h.Matches("correct-horse-battery-staple", hash) {
		t.Fatal("Matches rejected the correct password")
	}
}

func TestBcryptHasher_Matches_WrongPassword(t *testing.T) {
	h := NewBcryptHasher(4)

	hash, err := h.Hash("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}

	if h.Matches("wrong-password", hash) {
		t.Fatal("Matches accepted an incorrect password")
	}
}

func TestBcryptHasher_Matches_InvalidHash(t *testing.T) {
	h := NewBcryptHasher(4)

	if h.Matches("any-password", "not-a-valid-bcrypt-hash") {
		t.Fatal("Matches accepted a malformed hash")
	}
}
