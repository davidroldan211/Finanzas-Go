package domain

import "errors"

// ErrEmailRequired indica que la petición no trae email.
var ErrEmailRequired = errors.New("email is required")
