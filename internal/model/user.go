package model

import (
	"github.com/google/uuid"
)

type User struct {
	UserId    uuid.UUID
	FirstName string `validate:"required,min=2,max=50"`
	LastName  string `validate:"required,min=2,max=50"`
	Email     string `validate:"required,email"`
	Phone     string `validate:"required,e164"`
	Age       int    `validate:"omitempty,gt=0"`
	Status    Status `validate:"omitempty,oneof=Active Inactive"` // TODO: add default value
}

type Status string

const (
	StatusActive   Status = "Active"
	StatusInactive Status = "Inactive"
)
