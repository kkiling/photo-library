package session_manager

import (
	"crypto/rsa"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Claims .
type Claims interface {
	jwt.Claims
}

// NewHelper .
func NewHelper(
	cfg SslConfig,
) (JWTHelper, error) {
	buf, err := os.ReadFile(cfg.PrivateKeyFile)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(buf)
	if err != nil {
		return nil, err
	}

	buf, err = os.ReadFile(cfg.PublicKeyFile)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(buf)
	if err != nil {
		return nil, err
	}

	return &jwtHelper{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

type jwtHelper struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (h *jwtHelper) Parse(token string, claims Claims) error {
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return h.publicKey, nil
	})
	if err != nil {
		return err
	}

	return claims.Valid()
}

func (h *jwtHelper) CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(h.privateKey)
}
