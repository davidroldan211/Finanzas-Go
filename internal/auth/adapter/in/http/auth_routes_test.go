package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"finanzas-api/internal/auth/domain"
	"finanzas-api/internal/auth/port/in"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockAuthUseCase se define en auth_handler_test.go; ambos archivos
// comparten paquete tras la mudanza a adapter/in/http.

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockAuthUseCase{
		loginFunc: func(ctx context.Context, cmd in.LoginCommand) (domain.Token, error) {
			return domain.Token{Value: "mocked_token"}, nil
		},
	}

	router := gin.New()
	h := NewAuthHandler(mockUseCase)
	SetupAuthRoutes(router, h)

	w := httptest.NewRecorder()
	reqBody := `{"email":"test@example.com","password":"123456"}`
	req := httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mocked_token")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockAuthUseCase{
		loginFunc: func(ctx context.Context, cmd in.LoginCommand) (domain.Token, error) {
			return domain.Token{}, domain.ErrInvalidCredentials
		},
	}

	router := gin.New()
	h := NewAuthHandler(mockUseCase)
	SetupAuthRoutes(router, h)

	w := httptest.NewRecorder()
	reqBody := `{"email":"wrong@example.com","password":"wrongpass"}`
	req := httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"unauthorized"`)
}

func TestLogin_WrongMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockAuthUseCase{
		loginFunc: func(ctx context.Context, cmd in.LoginCommand) (domain.Token, error) {
			return domain.Token{Value: "mocked_token"}, nil
		},
	}

	router := gin.New()
	h := NewAuthHandler(mockUseCase)
	SetupAuthRoutes(router, h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/login", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
