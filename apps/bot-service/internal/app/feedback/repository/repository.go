package repository

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/jmoiron/sqlx"
)

type feedbackRepository struct {
	db *sqlx.DB
}

func NewFeedbackRepository(db *sqlx.DB) contracts.FeedbackRepository {
	return &feedbackRepository{db: db}
}
