package handler

import (
	"fmt"

	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/pkg/errors"
)

/*func mapFieldViolation(fieldViolations []serviceerr.FieldViolation) []*pbv1.FieldViolation {
	res := make([]*pbv1.FieldViolation, 0, len(fieldViolations))
	for _, f := range fieldViolations {
		res = append(res, &pbv1.FieldViolation{
			Field: f.Field,
			Error: f.ErrorMsg,
		})
	}
	return res
}*/

func handleError(err2 error, description any) error {
	newErr := fmt.Errorf("%s: %w", description, err2)

	/*info := pbv1.ErrorInfo{
		Description: "Unhandled error",
	}

	var serviceErr *serviceerr.ErrorService
	if errors.As(newErr, &serviceErr) {
		info = pbv1.ErrorInfo{
			Description:     serviceErr.ErrInfo.Description,
			FieldViolations: mapFieldViolation(serviceErr.ErrInfo.FieldViolations),
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
		// ErrPermissionDenied
		// ErrTooManyRequests
	}

	return server.ErrInternal(newErr, &info)*/

	switch {
	case errors.Is(newErr, serviceerr.ErrNotFound):
		return server.ErrNotFound(newErr)
	case errors.Is(newErr, serviceerr.ErrInvalidInput):
		return server.ErrInvalidArgument(newErr)
	case errors.Is(newErr, serviceerr.ErrConflict):
		return server.ErrAlreadyExists(newErr)
	}

	return server.ErrInternal(newErr)
}
