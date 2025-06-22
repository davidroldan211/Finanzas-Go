package routes

import (
	"errors"
	"finanzas-api/internal/auth/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAuthUseCase struct {
	loginFunc func(email, password string) (string, error)
}

func (m *mockAuthUseCase) Login(email, password string) (string, error) {
	return m.loginFunc(email, password)
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockAuthUseCase{
		loginFunc: func(email, password string) (string, error) {
			return "mocked_token", nil
		},
	}

	router := gin.New()
	h := handler.NewAuthHandler(mockUseCase)
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
		loginFunc: func(email, password string) (string, error) {
			return "", errors.New("invalid credentials")
		},
	}

	router := gin.New()
	h := handler.NewAuthHandler(mockUseCase)
	SetupAuthRoutes(router, h)

	w := httptest.NewRecorder()
	reqBody := `{"email":"wrong@example.com","password":"wrongpass"}`
	req := httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
}

func TestLogin_WrongMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockAuthUseCase{
		loginFunc: func(email, password string) (string, error) {
			return "mocked_token", nil
		},
	}

	router := gin.New()
	h := handler.NewAuthHandler(mockUseCase)
	SetupAuthRoutes(router, h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/login", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
