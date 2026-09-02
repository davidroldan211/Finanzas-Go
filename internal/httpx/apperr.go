package httpx

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`

	// Error interno (no se expone). Útil para logs y para wrap.
	Err error `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func New(status int, code, message string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func (e *AppError) WithErr(err error) *AppError {
	e.Err = err
	return e
}

func (e *AppError) WithFields(fields map[string]string) *AppError {
	if len(fields) == 0 {
		return e
	}
	if e.Fields == nil {
		e.Fields = map[string]string{}
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

func Is(err error) (*AppError, bool) {
	var app *AppError
	if errors.As(err, &app) {
		return app, true
	}
	return nil, false
}

// Wrap convierte cualquier error en AppError.
// Si ya es AppError, lo devuelve como está.
// Útil como última línea de defensa.
func Wrap(err error) *AppError {
	if err == nil {
		return nil
	}
	if app, ok := Is(err); ok {
		return app
	}
	return New(http.StatusInternalServerError, "internal_error", "Ocurrió un error inesperado.").WithErr(err)
}

/***************
 * Factories
 ***************/

func Validation(fields map[string]string) *AppError {
	return New(http.StatusUnprocessableEntity, "validation_error", "Hay campos inválidos.").WithFields(fields)
}

func BadRequest(msg string) *AppError {
	if msg == "" {
		msg = "Solicitud inválida."
	}
	return New(http.StatusBadRequest, "bad_request", msg)
}

func Unauthorized(msg string) *AppError {
	if msg == "" {
		msg = "No autorizado."
	}
	return New(http.StatusUnauthorized, "unauthorized", msg)
}

func Forbidden(msg string) *AppError {
	if msg == "" {
		msg = "No tienes permisos para realizar esta acción."
	}
	return New(http.StatusForbidden, "forbidden", msg)
}

func NotFound(msg string) *AppError {
	if msg == "" {
		msg = "Recurso no encontrado."
	}
	return New(http.StatusNotFound, "not_found", msg)
}

func Conflict(msg string) *AppError {
	if msg == "" {
		msg = "Conflicto con el estado actual del recurso."
	}
	return New(http.StatusConflict, "conflict", msg)
}

func TooManyRequests(msg string) *AppError {
	if msg == "" {
		msg = "Has hecho demasiadas solicitudes. Intenta más tarde."
	}
	return New(http.StatusTooManyRequests, "rate_limited", msg)
}
