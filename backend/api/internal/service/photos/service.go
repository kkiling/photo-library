package photosservice

import "github.com/kkiling/photo-library/backend/api/pkg/common/log"

type Service struct {
}

func NewService(logger log.Logger) *Service {
	return &Service{}
}
