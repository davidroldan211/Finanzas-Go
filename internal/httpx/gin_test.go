package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAbort_AppError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/x", nil)

	Abort(c, NotFound("Usuario no encontrado."))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"code":"not_found"`) {
		t.Fatalf("expected body to contain code not_found, got %s", body)
	}
	if !strings.Contains(body, "Usuario no encontrado.") {
		t.Fatalf("expected body to contain message, got %s", body)
	}
}

func TestAbort_PlainError_Returns500WithoutLeakingMessage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/x", nil)

	Abort(c, errors.New("conexión rechazada por la base de datos en 10.0.0.5"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"code":"internal_error"`) {
		t.Fatalf("expected body to contain code internal_error, got %s", body)
	}
	if strings.Contains(body, "10.0.0.5") {
		t.Fatalf("internal error message leaked to the client: %s", body)
	}
}

func TestAbortWith(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/x", nil)

	AbortWith(c, Conflict("El email ya está registrado."))

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

type bindTestPayload struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func TestBindJSON_ValidationError_PopulatesFields(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{"email":"not-an-email"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	var dst bindTestPayload
	ok := BindJSON(c, &dst)

	if ok {
		t.Fatal("expected BindJSON to return false on validation error")
	}
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"fields"`) {
		t.Fatalf("expected body to contain fields map, got %s", body)
	}
	if !strings.Contains(body, `"email"`) {
		t.Fatalf("expected fields to mention email, got %s", body)
	}
}

func TestBindJSON_MalformedBody_ReturnsBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`not-json`))
	c.Request.Header.Set("Content-Type", "application/json")

	var dst bindTestPayload
	ok := BindJSON(c, &dst)

	if ok {
		t.Fatal("expected BindJSON to return false on malformed body")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestBindJSON_ValidBody_ReturnsTrue(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{"email":"a@b.com","name":"A"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	var dst bindTestPayload
	ok := BindJSON(c, &dst)

	if !ok {
		t.Fatal("expected BindJSON to return true on valid body")
	}
	if dst.Email != "a@b.com" {
		t.Fatalf("expected email to be bound, got %q", dst.Email)
	}
}
