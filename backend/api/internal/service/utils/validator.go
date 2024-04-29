package utils

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func validateDurationMin(fl validator.FieldLevel) bool {
	param := fl.Param()
	minDuration, err := time.ParseDuration(param)
	if err != nil {
		return false
	}

	val, ok := fl.Field().Interface().(time.Duration)
	if !ok {
		return false
	}

	if val < minDuration {
		return false
	}

	return true
}

func validateDurationMax(fl validator.FieldLevel) bool {
	param := fl.Param()
	maxDuration, err := time.ParseDuration(param)
	if err != nil {
		return false
	}

	val, ok := fl.Field().Interface().(time.Duration)
	if !ok {
		return false
	}

	if val > maxDuration {
		return false
	}

	return true
}

func NewValidator() *validator.Validate {
	v := validator.New()
	if err := v.RegisterValidation("duration-min", validateDurationMin); err != nil {
		panic(err)
	}

	if err := v.RegisterValidation("duration-max", validateDurationMax); err != nil {
		panic(err)
	}
	return v
}
