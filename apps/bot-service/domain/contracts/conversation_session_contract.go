package contracts

import (
	"context"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=../../internal/app/session/repository/mock/mock_session_repository.go -package=mock github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts ConversationSessionRepository

type ConversationSessionRepository interface {
	Create(ctx context.Context, session *entity.ConversationSession) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ConversationSession, error)
	FindActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.ConversationSession, error)
	List(ctx context.Context, filter *entity.GetSessionsFilter) ([]entity.ConversationSession, int64, error)
	Update(ctx context.Context, session *entity.ConversationSession) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindExpiredSessions(ctx context.Context, inactiveMinutes int) ([]entity.ConversationSession, error)
}

type ConversationSessionService interface {
	GetOrCreateSession(ctx context.Context, phoneNumber string) (*entity.ConversationSession, error)
	UpdateSessionActivity(ctx context.Context, sessionID uuid.UUID) error
	MarkSessionWaitingFeedback(ctx context.Context, sessionID uuid.UUID) error
	CloseSession(ctx context.Context, sessionID uuid.UUID) error
	ProcessExpiredSessions(ctx context.Context) error
}
