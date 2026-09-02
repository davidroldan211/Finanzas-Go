package security

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"finanzas-api/internal/auth/domain"
)

func TestHMACTokenProvider_Issue_And_Verify_RoundTrip(t *testing.T) {
	p := NewHMACTokenProvider("test-secret", time.Hour)

	token, err := p.Issue(domain.Claims{Subject: "user-123", Role: "admin"})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}
	if token.Value == "" {
		t.Fatal("Issue returned an empty token")
	}

	claims, err := p.Verify(token.Value)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if claims.Subject != "user-123" || claims.Role != "admin" {
		t.Fatalf("expected Subject=user-123 Role=admin, got %+v", claims)
	}
}

func TestHMACTokenProvider_Verify_TamperedSignature(t *testing.T) {
	p := NewHMACTokenProvider("test-secret", time.Hour)

	token, err := p.Issue(domain.Claims{Subject: "user-123", Role: "user"})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	parts := strings.Split(token.Value, ".")
	tampered := parts[0] + "." + parts[1] + ".AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	if _, err := p.Verify(tampered); err == nil {
		t.Fatal("Verify accepted a token with a tampered signature")
	}
}

func TestHMACTokenProvider_Verify_WrongSecret(t *testing.T) {
	issuer := NewHMACTokenProvider("test-secret", time.Hour)
	verifier := NewHMACTokenProvider("different-secret", time.Hour)

	token, err := issuer.Issue(domain.Claims{Subject: "user-123", Role: "user"})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	if _, err := verifier.Verify(token.Value); err == nil {
		t.Fatal("Verify accepted a token signed with a different secret")
	}
}

func TestHMACTokenProvider_Verify_Expired(t *testing.T) {
	p := NewHMACTokenProvider("test-secret", -time.Hour)

	token, err := p.Issue(domain.Claims{Subject: "user-123", Role: "user"})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	if _, err := p.Verify(token.Value); err == nil {
		t.Fatal("Verify accepted an expired token")
	}
}

func TestHMACTokenProvider_Verify_MalformedToken(t *testing.T) {
	p := NewHMACTokenProvider("test-secret", time.Hour)

	if _, err := p.Verify("not-a-jwt"); err == nil {
		t.Fatal("Verify accepted a malformed token")
	}
}

func TestHMACTokenProvider_Verify_RejectsAlgNone(t *testing.T) {
	p := NewHMACTokenProvider("test-secret", time.Hour)

	header, _ := json.Marshal(tokenHeader{Alg: "none", Typ: "JWT"})
	payload, _ := json.Marshal(tokenPayload{UserID: "user-123", Role: "admin", Exp: time.Now().Add(time.Hour).Unix()})
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	forged := unsigned + "." + p.sign(unsigned)

	if _, err := p.Verify(forged); err == nil {
		t.Fatal("Verify accepted a token with alg=none")
	}
}

func TestHMACTokenProvider_BackwardCompatible_WithUUIDSubject(t *testing.T) {
	// Los tokens emitidos por la versión anterior llevaban un uuid.UUID en
	// user_id, que serializa a JSON como el mismo string que un Subject
	// de texto plano. Simula ese payload heredado y confirma que se
	// parsea igual.
	p := NewHMACTokenProvider("test-secret", time.Hour)

	header, _ := json.Marshal(tokenHeader{Alg: hmacAlg, Typ: "JWT"})
	legacyPayload := `{"user_id":"fc66d490-d765-4e04-bb09-bc9649902b75","role":"admin","exp":` +
		strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10) + `}`
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString([]byte(legacyPayload))
	token := unsigned + "." + p.sign(unsigned)

	claims, err := p.Verify(token)
	if err != nil {
		t.Fatalf("Verify returned error on legacy payload: %v", err)
	}
	if claims.Subject != "fc66d490-d765-4e04-bb09-bc9649902b75" {
		t.Fatalf("expected legacy uuid subject preserved, got %q", claims.Subject)
	}
}
