package errx

import (
	"net/http"
)

var (
	ErrFeedbackNotFound = NewError(
		http.StatusNotFound,
		"feedback_not_found",
		"Feedback not found.",
	)
	ErrFeedbackAlreadyExists = NewError(
		http.StatusConflict,
		"feedback_already_exists",
		"Feedback for this session already exists.",
	)
	ErrInvalidRating = NewError(
		http.StatusBadRequest,
		"invalid_rating",
		"Rating must be between 1 and 5.",
	)
)
