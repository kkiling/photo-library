package serviceerr

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrTagAlreadyExist = errors.New("tag already exist")
	ErrPhotoIsNotValid = errors.New("photo is not valid")
	ErrAlreadyLocked   = fmt.Errorf("rocket Key already locked")
)
