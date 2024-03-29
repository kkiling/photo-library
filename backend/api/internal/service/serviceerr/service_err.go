package serviceerr

import "errors"

var (
	ErrTagAlreadyExist = errors.New("tag already exist")
	ErrPhotoIsNotValid = errors.New("photo is not valid")
)
