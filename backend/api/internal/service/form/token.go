package form

import (
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"time"
)

type CreateApiToken struct {
	Caption      string             `validate:"required,min=3,max=128"`
	Type         model.ApiTokenType `validate:"required,min=3,max=128"`
	TimeDuration *time.Duration     `validate:"required,duration-min=24h,duration-max=4320h"`
}
