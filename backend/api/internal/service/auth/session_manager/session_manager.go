package session_manager

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
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
func (s *SessionManager) CreateTokenBySession(session model.Session) (Token, error) {
	expiresAt := time.Now().Add(s.cfg.AccessTokenDuration)
	accessClaims := &SessionClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  s.cfg.Audience,
			Issuer:    s.cfg.Issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
		Session: session,
	}
	accessToken, err := s.jwtHelper.CreateToken(accessClaims)
	if err != nil {
		return Token{}, serviceerr.MakeErr(err, "s.jwtHelper.CreateToken")
	}

	return Token{
		Token:     accessToken,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *SessionManager) CreateTokenByRefresh(refresh model.RefreshSession) (Token, error) {
	expiresAt := time.Now().Add(s.cfg.AccessTokenDuration)
	accessClaims := &RefreshSessionClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  s.cfg.Audience,
			Issuer:    s.cfg.Issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
		RefreshSession: refresh,
	}
	accessToken, err := s.jwtHelper.CreateToken(accessClaims)
	if err != nil {
		return Token{}, serviceerr.MakeErr(err, "s.jwtHelper.CreateToken")
	}

	return Token{
		Token:     accessToken,
		ExpiresAt: expiresAt,
	}, nil
}

// GetSessionByToken получить данные авторизованного пользователя по jwt токену
func (s *SessionManager) GetSessionByToken(token string) (model.Session, error) {
	claims := new(SessionClaims)
	err := s.jwtHelper.Parse(token, claims)
	if err != nil {
		return model.Session{}, server.ErrUnauthenticated(ErrSessionNotFound)
	}
	return claims.Session, nil
}

// GetRefreshSessionByToken получить данные авторизованного пользователя по jwt токену
func (s *SessionManager) GetRefreshSessionByToken(token string) (model.RefreshSession, error) {
	claims := new(RefreshSessionClaims)
	err := s.jwtHelper.Parse(token, claims)
	if err != nil {
		return model.RefreshSession{}, server.ErrUnauthenticated(ErrSessionNotFound)
	}
	return claims.RefreshSession, nil
}
