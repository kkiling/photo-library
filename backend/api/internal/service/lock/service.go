package lock

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

var RocketLockSafetyMargin = 500 * time.Millisecond

// Storage .
type Storage interface {
	RocketLock(ctx context.Context, key string, ttl time.Duration) (RocketLockID, error)
	RocketLockDelete(ctx context.Context, lockID RocketLockID) error
}

// Service .
type Service struct {
	storage Storage
}

// NewService новый сервис
func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func validateLockClientRequest(key string, ttl time.Duration) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("key must not be empty")
	}

	if ttl < 10*time.Millisecond {
		return fmt.Errorf("ttl must be greater than 10 milliseconds")
	}

	if ttl > 120*time.Second {
		return fmt.Errorf("ttl must be less than 3 seconds")
	}

	return nil
}

func (s *Service) Lock(ctx context.Context, key string, ttl time.Duration) (RocketLockID, error) {
	err := validateLockClientRequest(key, ttl)
	if err != nil {
		return RocketLockID{}, serviceerr.InvalidInputErr(err, "validateLockClientRequest")
	}

	lockID, err := s.storage.RocketLock(ctx, key, ttl)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrAlreadyLocked):
		return RocketLockID{}, err
	default:
		return RocketLockID{}, fmt.Errorf("s.storage.RocketLock: %w", err)
	}

	return lockID, nil
}

func (s *Service) UnLock(ctx context.Context, lockID RocketLockID) error {
	err := s.storage.RocketLockDelete(ctx, lockID)
	if err != nil {
		return fmt.Errorf("s.storage.RocketLockDelete: %w", err)
	}

	return nil
}
