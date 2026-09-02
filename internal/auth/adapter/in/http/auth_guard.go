package http

import (
	"log"
	nethttp "net/http"
	"strings"

	"finanzas-api/internal/shared/security"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Secret     string
	ParseToken ParseTokenFunc
}

type ParseTokenFunc func(token, secret string) (*security.TokenClaims, error)

func NewMiddleware(secret string, parseToken ParseTokenFunc) *Middleware {
	return &Middleware{
		Secret:     secret,
		ParseToken: parseToken,
	}
}

func (m *Middleware) Handler(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) > 0 && strings.EqualFold(fields[0], "bearer") {
			if len(fields) < 2 {
				c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			authHeader = fields[1]
		}

		authHeader = strings.TrimSpace(authHeader)
		claims, err := m.ParseToken(authHeader, m.Secret)
		if err != nil {
			c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
				c.AbortWithStatusJSON(nethttp.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		}
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		log.Printf("Authenticated user ID: %d with role: %s", claims.UserID, claims.Role)
		c.Next()
	}
}
