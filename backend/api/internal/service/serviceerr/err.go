package serviceerr

import (
	"fmt"
	"go.uber.org/multierr"

	"github.com/pkg/errors"
)

var ErrTagAlreadyExist = errors.New("tag already exist")
var ErrPhotoIsNotValid = errors.New("photo is not valid")

var ErrNotFound = errors.New("not found error")
var ErrInvalidInput = errors.New("invalid input")
var ErrConflict = errors.New("conflict err")

// InvalidInputError создание ErrInvalidInput
func InvalidInputError(description string, a ...any) error {
	err := fmt.Errorf(description, a...)
	return multierr.Append(err, ErrInvalidInput)
}

// InvalidInputErr создание ErrInvalidInput
func InvalidInputErr(err error, method string) error {
	err2 := fmt.Errorf(method+": %w", err)
	return multierr.Append(err2, ErrInvalidInput)
}

// NotFoundError создание NotFoundErrorType
func NotFoundError(description string, a ...any) error {
	err := fmt.Errorf(description, a...)
	return multierr.Append(err, ErrNotFound)
}

// ConflictError создание NotFoundErrorType
func ConflictError(description string, a ...any) error {
	err := fmt.Errorf(description, a...)
	return multierr.Append(err, ErrConflict)
}

// MakeErr создание RuntimeErrorType
func MakeErr(err error, method string) error {
	return fmt.Errorf(method+": %w", err)
}
