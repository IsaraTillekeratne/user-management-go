package mapper

import (
	"example.com/user-management/internal/dto"
	"example.com/user-management/internal/model"
)

func CreateUserRequestToModel(req dto.CreateUserRequest) model.User {
	return model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Age:       req.Age,
		Status:    req.Status,
	}
}

func ApplyUpdateUserRequest(u *model.User, req dto.UpdateUserRequest) {
	if req.FirstName != nil {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		u.LastName = *req.LastName
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Phone != nil {
		u.Phone = *req.Phone
	}
	if req.Age != nil {
		u.Age = *req.Age
	}
	if req.Status != nil {
		u.Status = *req.Status
	}
}
