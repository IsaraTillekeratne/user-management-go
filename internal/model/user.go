package model

import (
	"github.com/google/uuid"
)

type User struct {
	UserId    uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Age       int
	Status    Status
}

type Status string

const (
	StatusActive   Status = "Active"
	StatusInactive Status = "Inactive"
)
