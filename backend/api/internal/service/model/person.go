package model

import (
	"fmt"

	"github.com/google/uuid"
)

// Person информация о человеке
type Person struct {
	Base
	ID         uuid.UUID
	FirstName  string
	Surname    string
	Patronymic *string
}

// FullName полное имя человека
func (p Person) FullName() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.Surname)
}

type UpdatePerson struct {
	BaseUpdate
	FirstName  UpdateField[string]
	Surname    UpdateField[string]
	Patronymic UpdateField[*string]
}

// PersonDTO человек с полной информацией о нем
type PersonDTO struct {
	ID         uuid.UUID
	FirstName  string
	Surname    string
	Patronymic *string
	Email      string
}
