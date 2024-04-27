package session_manager

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/kkiling/photo-library/backend/api/internal/service/auth"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

var (
	// ErrSessionNotFound сессия авторизованного пользователя не найдена
	ErrSessionNotFound = errors.New("Session not found")
)

// JWTHelper хелпер для работы с jwt
type JWTHelper interface {
	Parse(token string, claims Claims) error
	CreateToken(claims Claims) (string, error)
}

// NewSessionManager ..
func NewSessionManager(
	cfg SessionConfig,
	logger log.Logger,
	jwtHelper JWTHelper,
) (*SessionManager, error) {
	return &SessionManager{
		logger:    logger.Named("session_manager"),
		cfg:       cfg,
		jwtHelper: jwtHelper,
	}, nil
}

// SessionManager менеджер работы с токенами и данными авторизованного пользователя
type SessionManager struct {
	logger    log.Logger
	cfg       SessionConfig
	jwtHelper JWTHelper
}

// CreateTokenBySession создание jwt токена для пользователя
func (s *SessionManager) CreateTokenBySession(_ context.Context, session auth.Session) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.cfg.AccessTokenDuration)
	claims := &SessionClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  s.cfg.Audience,
			Issuer:    s.cfg.Issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
		Session: session,
	}
	token, err := s.jwtHelper.CreateToken(claims)
	if err != nil {
		return "", time.Time{}, serviceerr.MakeErr(err, "s.jwtHelper.CreateToken")
	}
	return token, expiresAt, nil
}

// GetSessionByToken получить данные авторизованного пользователя по jwt токену
func (s *SessionManager) GetSessionByToken(ctx context.Context, token string) (auth.Session, error) {
	logger := s.logger.WithCtx(ctx)

	claims := new(SessionClaims)
	err := s.jwtHelper.Parse(token, claims)
	if err != nil {
		logger.Error(err)
		return auth.Session{}, server.ErrUnauthenticated(ErrSessionNotFound)
	}
	return claims.Session, nil
}
