package errx

import (
	"net/http"
)

var (
	ErrSessionNotFound = NewError(
		http.StatusNotFound,
		"session_not_found",
		"Conversation session not found.",
	)
	ErrSessionAlreadyClosed = NewError(
		http.StatusConflict,
		"session_already_closed",
		"This conversation session is already closed.",
	)
	ErrNoActiveSession = NewError(
		http.StatusNotFound,
		"no_active_session",
		"No active conversation session found for this user.",
	)
)
