package entity

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID          uuid.UUID `db:"id"`
	SessionID   uuid.UUID `db:"session_id"`
	PhoneNumber string    `db:"phone_number"`
	Rating      int       `db:"rating"`
	Comment     *string   `db:"comment"`
	CreatedAt   time.Time `db:"created_at"`
}

type FeedbackWithSession struct {
	Feedback
	Session ConversationSession
}

type GetFeedbacksFilter struct {
	Offset      int
	Limit       int
	PhoneNumber *string
	MinRating   *int
	MaxRating   *int
}
