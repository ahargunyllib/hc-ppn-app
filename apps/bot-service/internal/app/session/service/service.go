package service

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid"
)

type SessionService struct {
	sessionRepo contracts.ConversationSessionRepository
	uuidPkg     uuid.UUIDInterface
	logger      *log.Logger
}

func NewSessionService(
	sessionRepo contracts.ConversationSessionRepository,
	uuidService uuid.UUIDInterface,
	logger *log.Logger,
) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		uuidPkg:     uuidService,
		logger:      logger,
	}
}
