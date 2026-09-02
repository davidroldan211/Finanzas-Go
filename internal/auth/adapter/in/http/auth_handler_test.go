package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"finanzas-api/internal/auth/domain"
	"finanzas-api/internal/auth/port/in"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockAuthUseCase implementa in.AuthService. Se define una sola vez aquí;
// auth_routes_test.go (mismo paquete tras la mudanza a adapter/in/http) lo
// reutiliza.
type mockAuthUseCase struct {
	loginFunc func(ctx context.Context, cmd in.LoginCommand) (domain.Token, error)
}

func (m *mockAuthUseCase) Login(ctx context.Context, cmd in.LoginCommand) (domain.Token, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, cmd)
	}
	return domain.Token{}, nil
}

func setupRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", handler.Login)
	return r
}

func TestLoginSuccess(t *testing.T) {
	mockUseCase := &mockAuthUseCase{
		loginFunc: func(ctx context.Context, cmd in.LoginCommand) (domain.Token, error) {
			return domain.Token{Value: "mocked_token"}, nil
		},
	}
	h := NewAuthHandler(mockUseCase)
	routes := setupRouter(h)

	body := []byte(`{"email":"test@example.com","password":"password123"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mocked_token")
}

func TestLoginValidationError(t *testing.T) {
	mockUseCase := &mockAuthUseCase{}
	h := NewAuthHandler(mockUseCase)
	routes := setupRouter(h)

	body := []byte(`{"email":"bad-email", "password":""}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	routes.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), `"fields"`)
}

func TestLogin_Unauthorized(t *testing.T) {
	mockUC := &mockAuthUseCase{
		loginFunc: func(ctx context.Context, cmd in.LoginCommand) (domain.Token, error) {
			return domain.Token{}, domain.ErrInvalidCredentials
		},
	}
	h := NewAuthHandler(mockUC)
	router := setupRouter(h)

	body := []byte(`{"email":"test@example.com", "password":"wrong"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"unauthorized"`)
}

func TestLogin_GenericError_DoesNotLeakMessage(t *testing.T) {
	mockUC := &mockAuthUseCase{
		loginFunc: func(ctx context.Context, cmd in.LoginCommand) (domain.Token, error) {
			return domain.Token{}, errors.New("conexión perdida con postgres en 10.0.0.9")
		},
	}
	h := NewAuthHandler(mockUC)
	router := setupRouter(h)

	body := []byte(`{"email":"test@example.com", "password":"whatever"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NotContains(t, w.Body.String(), "10.0.0.9")
}
