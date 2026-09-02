package out

import "finanzas-api/internal/auth/domain"

// TokenProvider es el puerto de salida que emite y verifica tokens. El
// secreto y el TTL viven en el adaptador que lo implementa (p.ej.
// *security.HMACTokenProvider); ni la aplicación ni el dominio los conocen.
type TokenProvider interface {
	Issue(claims domain.Claims) (domain.Token, error)
	Verify(raw string) (domain.Claims, error)
}
