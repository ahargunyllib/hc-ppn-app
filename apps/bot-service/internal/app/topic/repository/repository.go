package repository

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts"
	"github.com/jmoiron/sqlx"
)

type topicRepository struct {
	db *sqlx.DB
}

func NewTopicRepository(db *sqlx.DB) contracts.TopicRepository {
	return &topicRepository{db: db}
}
