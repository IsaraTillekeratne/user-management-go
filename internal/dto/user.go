package dto

import "example.com/user-management/internal/model"

type CreateUserRequest struct {
	FirstName string       `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string       `json:"lastName" validate:"required,min=2,max=50"`
	Email     string       `json:"email" validate:"required,email"`
	Phone     string       `json:"phone" validate:"required,e164"`
	Age       int          `json:"age" validate:"omitempty,gt=0"`
	Status    model.Status `json:"status" validate:"omitempty,oneof=Active Inactive"`
}

type UpdateUserRequest struct {
	FirstName *string       `json:"firstName" validate:"omitempty,min=2,max=50"`
	LastName  *string       `json:"lastName" validate:"omitempty,min=2,max=50"`
	Email     *string       `json:"email" validate:"omitempty,email"`
	Phone     *string       `json:"phone" validate:"omitempty,e164"`
	Age       *int          `json:"age" validate:"omitempty,gt=0"`
	Status    *model.Status `json:"status" validate:"omitempty,oneof=Active Inactive"`
}
