package session_manager

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

// SessionClaims .
type SessionClaims struct {
	jwt.StandardClaims
	model.Session
}

// Valid валидность claims
func (c *SessionClaims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}
	return nil
}

// RefreshSessionClaims .
type RefreshSessionClaims struct {
	jwt.StandardClaims
	model.RefreshSession
}

// Valid валидность claims
func (c *RefreshSessionClaims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}

	return nil
}

type Token struct {
	Token     string
	ExpiresAt time.Time
}
