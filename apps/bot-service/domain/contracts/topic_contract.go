package contracts

import (
	"context"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
)

//go:generate mockgen -destination=../../internal/app/topic/repository/mock/mock_topic_repository.go -package=mock github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts TopicRepository

type TopicRepository interface {
	BulkCreate(ctx context.Context, topics []entity.Topic) error
	GetHotTopics(ctx context.Context) ([]entity.Topic, error)
	GetTopicsCount(ctx context.Context) (int, error)
}

type TopicService interface {
	BulkCreate(ctx context.Context, req *dto.BulkCreateTopicsRequest) error
	GetHotTopics(ctx context.Context) (*dto.GetHotTopicsResponse, error)
	GetTopicsCount(ctx context.Context) (*dto.GetTopicsCountResponse, error)
}
