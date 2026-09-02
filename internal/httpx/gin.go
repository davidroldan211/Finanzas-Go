package httpx

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthGuard es el puerto de entrada que exponen los módulos que protegen
// rutas. Permite que otros módulos registren rutas protegidas sin importar
// el módulo que las autentica.
type AuthGuard func(roles ...string) gin.HandlerFunc

// errorBody es el envelope único de error de la API.
type errorBody struct {
	Error *AppError `json:"error"`
}

// Abort escribe err como AppError JSON y aborta la cadena de Gin.
// Cualquier error no reconocido se convierte en 500 internal_error vía Wrap,
// sin filtrar el mensaje interno al cliente; ese mensaje solo se registra
// en el log del servidor.
func Abort(c *gin.Context, err error) {
	app := Wrap(err)
	if app.Status >= 500 && app.Err != nil {
		log.Printf("httpx: %s %s -> %v", c.Request.Method, c.Request.URL.Path, app.Err)
	}
	c.AbortWithStatusJSON(app.Status, errorBody{Error: app})
}

// AbortWith es el atajo para errores construidos in situ en el handler.
func AbortWith(c *gin.Context, app *AppError) {
	Abort(c, app)
}

// BindJSON hace ShouldBindJSON sobre dst y, si falla, aborta la petición con
// la respuesta apropiada: 422 con el mapa de campos si el fallo es de
// validación, o 400 si el cuerpo no es JSON válido. Devuelve false cuando ya
// respondió — el handler debe hacer return en ese caso.
func BindJSON(c *gin.Context, dst any) bool {
	err := c.ShouldBindJSON(dst)
	if err == nil {
		return true
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		Abort(c, Validation(fieldErrors(ve)))
		return false
	}

	Abort(c, BadRequest("El cuerpo de la solicitud no es JSON válido."))
	return false
}

// fieldErrors traduce validator.ValidationErrors a un mapa campo -> mensaje.
func fieldErrors(ve validator.ValidationErrors) map[string]string {
	fields := make(map[string]string, len(ve))
	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		fields[field] = fmt.Sprintf("no cumple la regla '%s'", fe.Tag())
	}
	return fields
}
