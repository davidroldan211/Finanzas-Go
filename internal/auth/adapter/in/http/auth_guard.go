package http

import (
	"strings"

	"finanzas-api/internal/auth/port/out"
	"finanzas-api/internal/httpx"

	"github.com/gin-gonic/gin"
)

// Guard es el adaptador de entrada que autentica y autoriza por rol.
// Autenticar/autorizar es el caso de uso de auth: solo cambia el adaptador
// que lo dispara (cabecera HTTP en vez de body de login), por eso vive
// aquí y no en un paquete "middleware" neutral.
type Guard struct {
	tokens out.TokenProvider
}

func NewGuard(tokens out.TokenProvider) *Guard {
	return &Guard{tokens: tokens}
}

// Handler satisface httpx.AuthGuard.
func (g *Guard) Handler(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			httpx.Abort(c, httpx.Unauthorized(""))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) > 0 && strings.EqualFold(fields[0], "bearer") {
			if len(fields) < 2 {
				httpx.Abort(c, httpx.Unauthorized(""))
				return
			}
			authHeader = fields[1]
		}

		authHeader = strings.TrimSpace(authHeader)
		claims, err := g.tokens.Verify(authHeader)
		if err != nil {
			httpx.Abort(c, httpx.Unauthorized(""))
			return
		}

		if len(roles) > 0 {
			allowed := false
			for _, r := range roles {
				if claims.Role == r {
					allowed = true
					break
				}
			}
			if !allowed {
				httpx.Abort(c, httpx.Forbidden(""))
				return
			}
		}

		c.Set("userID", claims.Subject)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}
