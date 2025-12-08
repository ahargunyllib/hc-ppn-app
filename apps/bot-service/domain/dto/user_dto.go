package dto

import (
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
)

type UserResponse struct {
	ID          string  `json:"id"`
	PhoneNumber string  `json:"phoneNumber"`
	Name        string  `json:"name"`
	JobTitle    *string `json:"jobTitle,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	DateOfBirth *string `json:"dateOfBirth,omitempty"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func ToUserResponse(user *entity.User) UserResponse {
	var dateOfBirth *string
	if user.DateOfBirth != nil {
		formatted := user.DateOfBirth.Format(time.DateOnly)
		dateOfBirth = &formatted
	}

	return UserResponse{
		ID:          user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
		JobTitle:    user.JobTitle,
		Gender:      user.Gender,
		DateOfBirth: dateOfBirth,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}
}

type CreateUserRequest struct {
	PhoneNumber string  `json:"phoneNumber" validate:"e164,required,min=10,max=20"`
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	JobTitle    *string `json:"jobTitle,omitempty" validate:"omitempty,max=255"`
	Gender      *string `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
	DateOfBirth *string `json:"dateOfBirth,omitempty"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type UpdateUserParam struct {
	ID string `param:"id" validate:"required,uuid"`
}

type UpdateUserRequest struct {
	PhoneNumber *string `json:"phoneNumber,omitempty" validate:"omitempty,e164,min=10,max=20"`
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	JobTitle    *string `json:"jobTitle,omitempty" validate:"omitempty,max=255"`
	Gender      *string `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
	DateOfBirth *string `json:"dateOfBirth,omitempty"`
}

type DeleteUserParam struct {
	ID string `param:"id" validate:"required,uuid"`
}

type GetUsersQuery struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty,max=255"`
}

type GetUsersResponse struct {
	Users []UserResponse `json:"users"`
	Meta  struct {
		Pagination PaginationResponse `json:"pagination"`
	} `json:"meta"`
}

type GetUserByIDParam struct {
	ID string `param:"id" validate:"required,uuid"`
}

type GetUserByIDResponse struct {
	User UserResponse `json:"user"`
}

type GetUserByPhoneNumberParam struct {
	PhoneNumber string `param:"phoneNumber" validate:"required,min=10,max=20"`
}

type GetUserByPhoneNumberResponse struct {
	User UserResponse `json:"user"`
}

type GetAllPhoneNumbersResponse struct {
	PhoneNumbers []string `json:"phoneNumbers"`
}

type GetUserMetricsResponse struct {
	TotalUsers int `json:"totalUsers"`
}

type ImportUsersFromCSVResponse struct {
	Total   int                       `json:"total"`
	Success int                       `json:"success"`
	Failed  int                       `json:"failed"`
	Errors  []ImportUsersFromCSVError `json:"errors"`
}

type ImportUsersFromCSVError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}
