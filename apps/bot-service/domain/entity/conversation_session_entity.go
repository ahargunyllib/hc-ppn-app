package entity

import (
	"time"

	"github.com/google/uuid"
)

type ConversationSessionStatus string

const (
	SessionStatusActive          ConversationSessionStatus = "active"
	SessionStatusWaitingFeedback ConversationSessionStatus = "waiting_feedback"
	SessionStatusClosed          ConversationSessionStatus = "closed"
)

type ConversationSession struct {
	ID                   uuid.UUID                 `db:"id"`
	PhoneNumber          string                    `db:"phone_number"`
	Status               ConversationSessionStatus `db:"status"`
	LastMessageAt        time.Time                 `db:"last_message_at"`
	FeedbackPromptSentAt *time.Time                `db:"feedback_prompt_sent_at"`
	CreatedAt            time.Time                 `db:"created_at"`
	UpdatedAt            time.Time                 `db:"updated_at"`
}

type GetSessionsFilter struct {
	Offset      int
	Limit       int
	PhoneNumber *string
	Status      *ConversationSessionStatus
}
