package security

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateToken_And_ParseToken_RoundTrip(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := GenerateToken(userID, "admin", secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("expected UserID %v, got %v", userID, claims.UserID)
	}
	if claims.Role != "admin" {
		t.Fatalf("expected Role %q, got %q", "admin", claims.Role)
	}
}

func TestParseToken_TamperedSignature(t *testing.T) {
	token, err := GenerateToken(uuid.New(), "user", "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 token parts, got %d", len(parts))
	}
	tampered := parts[0] + "." + parts[1] + "." + "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	if _, err := ParseToken(tampered, "test-secret"); err == nil {
		t.Fatal("ParseToken accepted a token with a tampered signature")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken(uuid.New(), "user", "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if _, err := ParseToken(token, "different-secret"); err == nil {
		t.Fatal("ParseToken accepted a token signed with a different secret")
	}
}

func TestParseToken_Expired(t *testing.T) {
	token, err := GenerateToken(uuid.New(), "user", "test-secret", -time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if _, err := ParseToken(token, "test-secret"); err == nil {
		t.Fatal("ParseToken accepted an expired token")
	}
}

func TestParseToken_MalformedToken(t *testing.T) {
	if _, err := ParseToken("not-a-jwt", "test-secret"); err == nil {
		t.Fatal("ParseToken accepted a malformed token")
	}
}
