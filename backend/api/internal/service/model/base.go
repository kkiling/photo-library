package model

import "time"

// Base .
type Base struct {
	CreateAt time.Time
	UpdateAt time.Time
}

// BaseUpdate .
type BaseUpdate struct {
	UpdateAt time.Time
}

// NewBase .
func NewBase() Base {
	return Base{
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
}

// NewBaseUpdate .
func NewBaseUpdate() BaseUpdate {
	return BaseUpdate{
		UpdateAt: time.Now(),
	}
}

// UpdateField .
type UpdateField[T any] struct {
	NeedUpdate bool
	Value      T
}

// NewUpdateField .
func NewUpdateField[T any](value T) UpdateField[T] {
	return UpdateField[T]{
		NeedUpdate: true,
		Value:      value,
	}
}
