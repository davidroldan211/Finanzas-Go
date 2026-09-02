package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"finanzas-api/internal/auth/domain"
)

const hmacAlg = "HS256"

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// tokenPayload es el payload serializado del token. Se mantiene el tag
// "user_id" (heredado de la versión anterior, donde el campo era
// uuid.UUID) para que los tokens emitidos antes de esta migración sigan
// siendo válidos: la representación JSON de un uuid.UUID ya era un string.
type tokenPayload struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

// HMACTokenProvider implementa auth/port/out.TokenProvider con HMAC-SHA256
// hecho a mano (sin librería JWT externa). Guarda el secreto y el TTL: el
// core de auth nunca los conoce.
type HMACTokenProvider struct {
	secret string
	ttl    time.Duration
}

func NewHMACTokenProvider(secret string, ttl time.Duration) *HMACTokenProvider {
	return &HMACTokenProvider{secret: secret, ttl: ttl}
}

func (p *HMACTokenProvider) Issue(claims domain.Claims) (domain.Token, error) {
	expiresAt := time.Now().Add(p.ttl)

	headerBytes, err := json.Marshal(tokenHeader{Alg: hmacAlg, Typ: "JWT"})
	if err != nil {
		return domain.Token{}, err
	}
	payloadBytes, err := json.Marshal(tokenPayload{
		UserID: claims.Subject,
		Role:   claims.Role,
		Exp:    expiresAt.Unix(),
	})
	if err != nil {
		return domain.Token{}, err
	}

	header := base64.RawURLEncoding.EncodeToString(headerBytes)
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	unsigned := header + "." + payload

	return domain.Token{
		Value:     unsigned + "." + p.sign(unsigned),
		ExpiresAt: expiresAt,
	}, nil
}

func (p *HMACTokenProvider) Verify(raw string) (domain.Claims, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 3 {
		return domain.Claims{}, domain.ErrInvalidToken
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return domain.Claims{}, domain.ErrInvalidToken
	}
	var header tokenHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return domain.Claims{}, domain.ErrInvalidToken
	}
	if header.Alg != hmacAlg {
		return domain.Claims{}, domain.ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(p.sign(unsigned)), []byte(parts[2])) {
		return domain.Claims{}, domain.ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return domain.Claims{}, domain.ErrInvalidToken
	}
	var payload tokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return domain.Claims{}, domain.ErrInvalidToken
	}

	expiresAt := time.Unix(payload.Exp, 0)
	if time.Now().After(expiresAt) {
		return domain.Claims{}, domain.ErrInvalidToken
	}

	return domain.Claims{Subject: payload.UserID, Role: payload.Role, ExpiresAt: expiresAt}, nil
}

func (p *HMACTokenProvider) sign(unsigned string) string {
	mac := hmac.New(sha256.New, []byte(p.secret))
	mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
