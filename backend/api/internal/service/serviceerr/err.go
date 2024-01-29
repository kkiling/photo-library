package serviceerr

import (
	"fmt"
	"github.com/kkiling/photo-library/backend/api/pkg/common/utils"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// ErrorServiceType типы ошибок сервисов
type ErrorServiceType int

// .
const (
	InvalidInputDataErrorType ErrorServiceType = 0
	RuntimeErrorType          ErrorServiceType = 1
	NotFoundErrorType         ErrorServiceType = 2
	ConflictErrorType         ErrorServiceType = 3
)

type FieldViolation struct {
	Field    string
	ErrorMsg string
}

type ErrorInfo struct {
	Description     string
	FieldViolations []FieldViolation
}

// ErrorService ошибка сервиса
type ErrorService struct {
	Type    ErrorServiceType
	Err     error
	ErrInfo ErrorInfo
}

func (r *ErrorService) Error() string {
	if r == nil || r.Err == nil {
		return ""
	}
	return r.Err.Error()
}

// InvalidInputValidatorError создание InvalidInputDataErrorType
func InvalidInputValidatorError(validatorError error, description string) error {
	if description == "" {
		description = "Data validation error"
	}
	var ve validator.ValidationErrors
	fv := make([]FieldViolation, 0)
	if errors.As(validatorError, &ve) {
		for _, fe := range ve {
			fv = append(fv, FieldViolation{
				Field:    fe.Field(),
				ErrorMsg: fmt.Sprintf("failed on the '%s' tag", fe.Tag()),
			})
		}
	}

	// , description string
	return &ErrorService{
		Type: InvalidInputDataErrorType,
		Err:  fmt.Errorf("fail validate data"+" :%w", validatorError),
		ErrInfo: ErrorInfo{
			Description:     description,
			FieldViolations: fv,
		},
	}
}

// InvalidInputError создание NotFoundErrorType
func InvalidInputError(description string, a ...any) error {
	err := fmt.Errorf(description, a...)
	return &ErrorService{
		Type: NotFoundErrorType,
		Err:  err,
		ErrInfo: ErrorInfo{
			Description: err.Error(),
		},
	}
}

// NotFoundError создание NotFoundErrorType
func NotFoundError(description string, a ...any) error {
	err := fmt.Errorf(description, a...)
	return &ErrorService{
		Type: NotFoundErrorType,
		Err:  err,
		ErrInfo: ErrorInfo{
			Description: err.Error(),
		},
	}
}

// ConflictError создание ConflictErrorType
func ConflictError(description string, a ...any) error {
	err := fmt.Errorf(description, a...)
	return &ErrorService{
		Type: ConflictErrorType,
		Err:  err,
		ErrInfo: ErrorInfo{
			Description: err.Error(),
		},
	}
}

// RuntimeError создание RuntimeErrorType
func RuntimeError(err error, method any) error {
	return &ErrorService{
		Err:  fmt.Errorf(utils.GetFunctionName(method)+": %w", err),
		Type: RuntimeErrorType,
		ErrInfo: ErrorInfo{
			Description: "Unknown runtime error",
		},
	}
}
