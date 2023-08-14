package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func validateInputs(name, color string) error {
	validate := validator.New()

	// Валидация имени
	if err := validate.Var(name, "min=3"); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return validationErrors
		}

		return fmt.Errorf("invalid name: %w", err)
	}

	// Валидация цвета в формате HEX (например, "#FFFFFF")
	// Вы можете настроить этот шаблон, если у вас есть другие требования к формату.
	if err := validate.Var(color, "hexcolor"); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return validationErrors
		}

		return fmt.Errorf("invalid color: %w", err)
	}

	return nil
}

func main() {
	name := "Al"
	color := "#G12345" // это некорректный HEX-цвет

	if err := validateInputs(name, color); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Validation successful!")
	}
}
