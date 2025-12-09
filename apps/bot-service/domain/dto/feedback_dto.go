package dto

import (
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
)

type FeedbackResponse struct {
	ID        string       `json:"id"`
	User      UserResponse `json:"user"`
	Rating    int          `json:"rating"`
	Comment   *string      `json:"comment,omitempty"`
	CreatedAt string       `json:"createdAt"`
}

func ToFeedbackResponse(feedback *entity.Feedback) FeedbackResponse {
	if feedback == nil {
		return FeedbackResponse{}
	}
	
	return FeedbackResponse{
		ID:        feedback.ID.String(),
		User:      ToUserResponse(&feedback.User),
		Rating:    feedback.Rating,
		Comment:   feedback.Comment,
		CreatedAt: feedback.CreatedAt.Format(time.RFC3339),
	}
}

type CreateFeedbackRequest struct {
	UserID  string  `json:"userId" validate:"required,uuid"`
	Rating  int     `json:"rating" validate:"required,min=1,max=5"`
	Comment *string `json:"comment,omitempty" validate:"omitempty,max=1000"`
}

type CreateFeedbackResponse struct {
	ID string `json:"id"`
}

type GetFeedbacksQuery struct {
	Page      int     `query:"page" validate:"omitempty,min=1"`
	Limit     int     `query:"limit" validate:"omitempty,min=1,max=100"`
	UserID    *string `query:"userId" validate:"omitempty,uuid"`
	Ratings   []int   `query:"ratings" validate:"omitempty,dive,min=1,max=5"`
	MinRating *int    `query:"minRating" validate:"omitempty,min=1,max=5"`
	MaxRating *int    `query:"maxRating" validate:"omitempty,min=1,max=5"`
}

type GetFeedbacksResponse struct {
	Feedbacks []FeedbackResponse `json:"feedbacks"`
	Meta      struct {
		Pagination PaginationResponse `json:"pagination"`
	} `json:"meta"`
}

type GetFeedbackByIDParam struct {
	ID string `param:"id" validate:"required,uuid"`
}

type GetFeedbackByIDResponse struct {
	Feedback FeedbackResponse `json:"feedback"`
}

type GetFeedbackMetricsResponse struct {
	SatisfactionScore float64 `json:"satisfactionScore"`
	TotalFeedbacks    int     `json:"totalFeedbacks"`
}

type SatisfactionTrendData struct {
	Date            string  `json:"date"`
	AvgSatisfaction float64 `json:"avgSatisfaction"`
}

type GetSatisfactionTrendResponse struct {
	Trend []SatisfactionTrendData `json:"trend"`
}
