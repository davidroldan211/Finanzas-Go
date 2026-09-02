package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAuthUseCase struct {
	loginFunc func(email, password string) (string, error)
}

func (m *mockAuthUseCase) Login(email, password string) (string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(email, password)
	}
	return "", nil
}

func setupRouter(handler *AuthHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/login", handler.Login)
	return r
}

func TestLoginSuccess(t *testing.T) {
	mockUseCase := &mockAuthUseCase{
		loginFunc: func(email, password string) (string, error) {
			return "mocked_token", nil
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

func TestLoginBadRequest(t *testing.T) {
	mockUseCase := &mockAuthUseCase{}
	h := NewAuthHandler(mockUseCase)
	routes := setupRouter(h)

	body := []byte(`{"email":"bad-email", "password":""}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	routes.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestLogin_Unauthorized(t *testing.T) {
	mockUC := &mockAuthUseCase{
		loginFunc: func(email, password string) (string, error) {
			return "", errors.New("invalid credentials")
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
	assert.Contains(t, w.Body.String(), "invalid credentials")
}
