package contracts

import (
	"context"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=../../internal/app/feedback/repository/mock/mock_feedback_repository.go -package=mock github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts FeedbackRepository

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *entity.Feedback) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Feedback, error)
	List(ctx context.Context, filter *entity.GetFeedbacksFilter) ([]entity.Feedback, int64, error)
	GetMetrics(ctx context.Context) (float64, int, error)
	GetSatisfactionTrend(ctx context.Context) ([]entity.SatisfactionTrendRow, error)
}

type FeedbackService interface {
	Create(ctx context.Context, req *dto.CreateFeedbackRequest) (*dto.CreateFeedbackResponse, error)
	GetByID(ctx context.Context, param *dto.GetFeedbackByIDParam) (*dto.GetFeedbackByIDResponse, error)
	List(ctx context.Context, query *dto.GetFeedbacksQuery) (*dto.GetFeedbacksResponse, error)
	GetMetrics(ctx context.Context) (*dto.GetFeedbackMetricsResponse, error)
	GetSatisfactionTrend(ctx context.Context) (*dto.GetSatisfactionTrendResponse, error)
}
