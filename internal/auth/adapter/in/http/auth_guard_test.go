package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"finanzas-api/internal/auth/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type stubTokenProvider struct {
	verifyFunc func(raw string) (domain.Claims, error)
}

func (s *stubTokenProvider) Issue(claims domain.Claims) (domain.Token, error) {
	return domain.Token{}, nil
}

func (s *stubTokenProvider) Verify(raw string) (domain.Claims, error) {
	return s.verifyFunc(raw)
}

func adminTokenProvider() *stubTokenProvider {
	return &stubTokenProvider{verifyFunc: func(raw string) (domain.Claims, error) {
		return domain.Claims{Subject: "12345678-1234-1234-1234-123456789012", Role: "admin"}, nil
	}}
}

func userTokenProvider() *stubTokenProvider {
	return &stubTokenProvider{verifyFunc: func(raw string) (domain.Claims, error) {
		return domain.Claims{Subject: "12345678-1234-1234-1234-123456789012", Role: "user"}, nil
	}}
}

func invalidTokenProvider() *stubTokenProvider {
	return &stubTokenProvider{verifyFunc: func(raw string) (domain.Claims, error) {
		return domain.Claims{}, errors.New("invalid token")
	}}
}

func expiredTokenProvider() *stubTokenProvider {
	return &stubTokenProvider{verifyFunc: func(raw string) (domain.Claims, error) {
		return domain.Claims{}, domain.ErrInvalidToken
	}}
}

func newProtectedRouter(g *Guard, roles ...string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", g.Handler(roles...), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	return r
}

func TestGuard_Unauthorized_InvalidToken(t *testing.T) {
	g := NewGuard(invalidTokenProvider())
	r := newProtectedRouter(g, "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
}

func TestGuard_Unauthorized_ExpiredToken(t *testing.T) {
	g := NewGuard(expiredTokenProvider())
	r := newProtectedRouter(g, "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer expiredtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
}

func TestGuard_Forbidden_WrongRole(t *testing.T) {
	g := NewGuard(userTokenProvider())
	r := newProtectedRouter(g, "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "forbidden")
}

func TestGuard_Success(t *testing.T) {
	g := NewGuard(adminTokenProvider())
	r := newProtectedRouter(g, "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestGuard_SuccessWithExtraSpaces(t *testing.T) {
	g := NewGuard(adminTokenProvider())
	r := newProtectedRouter(g, "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "  Bearer    validtoken   ")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestGuard_Unauthorized_NoHeader(t *testing.T) {
	g := NewGuard(adminTokenProvider())
	r := newProtectedRouter(g, "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGuard_Unauthorized_HeaderWithoutBearer_InvalidToken(t *testing.T) {
	g := NewGuard(invalidTokenProvider())
	r := newProtectedRouter(g, "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "not-a-real-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Sin el prefijo "Bearer", el header completo se trata como token y
	// pasa a Verify tal cual; con un token inválido, Verify lo rechaza -> 401.
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGuard_Unauthorized_BearerWithoutToken(t *testing.T) {
	g := NewGuard(adminTokenProvider())
	r := newProtectedRouter(g, "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
