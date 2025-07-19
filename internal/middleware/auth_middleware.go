package middleware

import (
	"net/http"
	"strings"

	"finanzas-api/shared/security"

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
		authHeader := c.GetHeader("Authorization")
		if after, ok := strings.CutPrefix(authHeader, "Bearer"); ok {
			authHeader = after
		}
		claims, err := m.ParseToken(authHeader, m.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		}
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}
