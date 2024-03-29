package handler

import (
	"fmt"

	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/pkg/errors"
)

func HandleError(err error, description any) error {
	newErr := fmt.Errorf("%s: %w", description, err)

	info := pbv1.ErrorInfo{
		Description: "Unhandled error",
	}

	var serviceErr *serviceerr.ErrorService
	if errors.As(newErr, &serviceErr) {
		info = pbv1.ErrorInfo{
			Description: serviceErr.ErrInfo.Description,
			// FieldViolations: mapFieldViolation(serviceErr.ErrInfo.FieldViolations),
		}
		switch serviceErr.Type {
		case serviceerr.InvalidInputDataErrorType:
			return server.ErrInvalidArgument(newErr, &info)
		case serviceerr.RuntimeErrorType:
			return server.ErrInternal(newErr, &info)
		case serviceerr.NotFoundErrorType:
			return server.ErrNotFound(newErr, &info)
		case serviceerr.ConflictErrorType:
			return server.ErrAlreadyExists(newErr, &info)
		}
	}

	return server.ErrInternal(newErr, &info)
}
