package repository

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) contracts.ConversationSessionRepository {
	return &sessionRepository{db: db}
}
