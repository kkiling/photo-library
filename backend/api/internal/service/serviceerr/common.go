package serviceerr

import (
	"errors"
	"fmt"
)

// ErrorServiceType типы ошибок сервисов
type ErrorServiceType int

const (
	InvalidInputDataErrorType ErrorServiceType = 0
	RuntimeErrorType          ErrorServiceType = 1
	NotFoundErrorType         ErrorServiceType = 2
	ConflictErrorType         ErrorServiceType = 3
	FailPreconditionErrorType ErrorServiceType = 4
	PermissionDeniedType      ErrorServiceType = 5
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

func (r *ErrorService) Unwrap() error {
	return r.Err
}

func IsNotFound(err error) bool {
	var serviceErr *ErrorService
	if errors.As(err, &serviceErr) {
		return serviceErr.Type == NotFoundErrorType
	}
	return false
}

// InvalidInputf создание ошибки входных данных
func InvalidInputf(description string, a ...any) error {
	return &ErrorService{
		Type: InvalidInputDataErrorType,
		Err:  fmt.Errorf(description, a...),
		ErrInfo: ErrorInfo{
			Description: "Invalid input",
		},
	}
}

// InvalidInputErr создание ошибки входных данных
func InvalidInputErr(err error, method string) error {
	return &ErrorService{
		Type: InvalidInputDataErrorType,
		Err:  fmt.Errorf(method+": %w", err),
		ErrInfo: ErrorInfo{
			Description: "Invalid input",
		},
	}
}

// NotFoundf создание ошибки - не найдено
func NotFoundf(description string, a ...any) error {
	return &ErrorService{
		Type: NotFoundErrorType,
		Err:  fmt.Errorf(description, a...),
		ErrInfo: ErrorInfo{
			Description: "Not found",
		},
	}
}

// Conflictf создание ошибки конфликта
func Conflictf(description string, a ...any) error {
	return &ErrorService{
		Type: ConflictErrorType,
		Err:  fmt.Errorf(description, a...),
		ErrInfo: ErrorInfo{
			Description: "Conflict",
		},
	}
}

// FailPreconditionf создание ошибки fail prediction
func FailPreconditionf(description string, a ...any) error {
	return &ErrorService{
		Type: FailPreconditionErrorType,
		Err:  fmt.Errorf(description, a...),
		ErrInfo: ErrorInfo{
			Description: "FailPrecondition",
		},
	}
}

func PermissionDeniedf(description string, a ...any) error {
	return &ErrorService{
		Type: PermissionDeniedType,
		Err:  fmt.Errorf(description, a...),
		ErrInfo: ErrorInfo{
			Description: "PermissionDenied",
		},
	}
}

// PermissionDeniedErr .
func PermissionDeniedErr(err error) error {
	return &ErrorService{
		Type: PermissionDeniedType,
		Err:  err,
		ErrInfo: ErrorInfo{
			Description: "Permission Denied",
		},
	}
}

// MakeErr создание
func MakeErr(err error, method string) error {
	return &ErrorService{
		Type: RuntimeErrorType,
		Err:  fmt.Errorf(method+": %w", err),
		ErrInfo: ErrorInfo{
			Description: "Runtime error",
		},
	}
}
