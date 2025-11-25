package errx

import (
	"net/http"
)

var (
	ErrUserNotFound = NewError(
		http.StatusNotFound,
		"user_not_found",
		"User not found.",
	)
	ErrUserPhoneExists = NewError(
		http.StatusConflict,
		"user_phone_exists",
		"A user with this phone number already exists.",
	)
	ErrUserInvalidPhone = NewError(
		http.StatusBadRequest,
		"user_invalid_phone",
		"Invalid phone number format.",
	)
)
