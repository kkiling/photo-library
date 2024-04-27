package session_manager

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/kkiling/photo-library/backend/api/internal/service/auth"
)

// SessionClaims .
type SessionClaims struct {
	jwt.StandardClaims
	auth.Session
}

// Valid валидность claims
func (c *SessionClaims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}

	return nil
}
