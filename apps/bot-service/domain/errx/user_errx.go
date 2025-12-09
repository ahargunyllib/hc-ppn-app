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
	ErrEmptyCSVFile = NewError(
		http.StatusBadRequest,
		"empty_csv_file",
		"CSV file is empty.",
	)
	ErrCSVNoData = NewError(
		http.StatusBadRequest,
		"csv_no_data",
		"CSV file must contain at least a header row and one data row.",
	)
	ErrInvalidCSVStructure = NewError(
		http.StatusBadRequest,
		"invalid_csv_structure",
		"CSV must have exactly 5 columns: phone_number, name, job_title, gender, date_of_birth.",
	)
	ErrInvalidCSVRow = NewError(
		http.StatusBadRequest,
		"invalid_csv_row",
		"CSV row has incorrect number of columns.",
	)
	ErrMissingPhoneNumber = NewError(
		http.StatusBadRequest,
		"missing_phone_number",
		"Phone number is required.",
	)
	ErrMissingName = NewError(
		http.StatusBadRequest,
		"missing_name",
		"Name is required.",
	)
	ErrInvalidPhoneNumber = NewError(
		http.StatusBadRequest,
		"invalid_phone_number",
		"Invalid phone number format (must be E.164 format).",
	)
	ErrInvalidGender = NewError(
		http.StatusBadRequest,
		"invalid_gender",
		"Gender must be either 'male' or 'female'.",
	)
)
