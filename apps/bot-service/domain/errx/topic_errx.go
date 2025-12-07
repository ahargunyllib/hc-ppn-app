package errx

import "net/http"

var (
	ErrTopicNotFound = NewError(http.StatusNotFound, "TOPIC_NOT_FOUND", "topic not found")
)
