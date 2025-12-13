package middleware

import (
	"errors"
	"finanzas-api/shared/security"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock: token válido con rol admin
func mockParseTokenAdmin(token, secret string) (*security.TokenClaims, error) {
	return &security.TokenClaims{
		UserID: 123,
		Role:   "admin",
		Exp:    9999999999,
	}, nil
}

// Mock: token válido con rol user
func mockParseTokenUser(token, secret string) (*security.TokenClaims, error) {
	return &security.TokenClaims{
		UserID: 123,
		Role:   "user",
		Exp:    9999999999,
	}, nil
}

// Mock: token inválido
func mockParseTokenInvalid(token, secret string) (*security.TokenClaims, error) {
	return nil, errors.New("invalid token")
}

// Mock: token expirado
func mockParseTokenExpired(token, secret string) (*security.TokenClaims, error) {
	return nil, errors.New("token expired")
}

func TestMiddleware_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	m := NewMiddleware("test-secret", mockParseTokenInvalid)
	r.GET("/protected", m.Handler("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
}

func TestMiddleware_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	m := NewMiddleware("test-secret", mockParseTokenUser)
	r.GET("/protected", m.Handler("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "forbidden")
}

func TestMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	m := NewMiddleware("test-secret", mockParseTokenAdmin)
	r.GET("/protected", m.Handler("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestMiddleware_SuccessWithExtraSpaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	m := NewMiddleware("test-secret", mockParseTokenAdmin)
	r.GET("/protected", m.Handler("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "  Bearer    validtoken   ")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}
