package errx

import (
	"net/http"
)

type RequestError struct {
	StatusCode int            `json:"-"`
	ErrorCode  string         `json:"error_code"`
	Message    string         `json:"message"`
	Location   string         `json:"-"`
	Details    map[string]any `json:"-"`
	Err        error          `json:"-"`
}

func (r *RequestError) Error() string {
	return r.Message
}

func NewError(statusCode int, errorCode string, message string) *RequestError {
	return &RequestError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (r *RequestError) WithLocation(location string) *RequestError {
	r.Location = location
	return r
}

func (r *RequestError) WithDetails(details map[string]any) *RequestError {
	r.Details = details
	return r
}

func (r *RequestError) WithError(err error) *RequestError {
	r.Err = err
	return r
}

func (r *RequestError) Unwrap() error {
	return r.Err
}

var (
	ErrInternalServer = NewError(
		http.StatusInternalServerError,
		"internal_server_error",
		"Something went wrong on our end. Please try again later.",
	)
	ErrNotFound = NewError(
		http.StatusNotFound,
		"not_found",
		"The requested resource could not be found.",
	)
	ErrForbidden = NewError(
		http.StatusForbidden,
		"forbidden",
		"You don't have permission to access this resource.",
	)
	ErrUnauthorized = NewError(
		http.StatusUnauthorized,
		"unauthorized",
		"You are not authorized to access this resource.",
	)

	ErrNoAPIKey = NewError(
		http.StatusUnauthorized,
		"no_api_key",
		"An API key is required to access this service.",
	)
	ErrInvalidAPIKey = NewError(
		http.StatusUnauthorized,
		"invalid_api_key",
		"The provided API key is not valid.",
	)
	ErrNoBearerToken = NewError(
		http.StatusUnauthorized,
		"no_bearer_token",
		"A bearer token is required to continue.",
	)
	ErrInvalidBearerToken = NewError(
		http.StatusUnauthorized,
		"invalid_bearer_token",
		"The bearer token provided is invalid.",
	)
	ErrExpiredBearerToken = NewError(
		http.StatusUnauthorized,
		"expired_bearer_token",
		"Your session has expired. Please log in again.",
	)
	ErrBearerTokenNotActive = NewError(
		http.StatusUnauthorized,
		"bearer_token_not_active",
		"Your token is not yet active. Please check and try again.",
	)
	ErrRoleCantAccessResource = NewError(
		http.StatusForbidden,
		"role_cant_access_resource",
		"You don't have permission to access this resource.",
	)
	ErrClaimsNotFound = NewError(
		http.StatusUnauthorized,
		"claims_not_found",
		"Authorization details are missing. Please try again.",
	)
	ErrInvalidTimeFormat = NewError(
		http.StatusBadRequest,
		"invalid_time_format",
		"Please enter a valid time format.",
	)
	ErrFileSizeLimitExceeded = NewError(
		http.StatusBadRequest,
		"file_size_limit_exceeded",
		"File size limit exceeded. Please try again with a smaller file.",
	)
	ErrInvalidFileExtension = NewError(
		http.StatusBadRequest,
		"invalid_file_extension",
		"Invalid file extension. Please try again with a different file.",
	)
	ErrTooManyRequests = NewError(
		http.StatusTooManyRequests,
		"too_many_requests",
		"Too many requests. Please try again later.",
	)
	ErrCannotParseUUID = NewError(
		http.StatusBadRequest,
		"cannot_parse_uuid",
		"Cannot parse the provided UUID. Please check the format and try again.",
	)
	ErrInvalidDateFormat = NewError(
		http.StatusBadRequest,
		"invalid_date_format",
		"Please enter a valid date format.",
	)
)
