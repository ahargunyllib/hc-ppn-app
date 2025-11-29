package dto

import (
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
)

type FeedbackResponse struct {
	ID          string  `json:"id"`
	SessionID   string  `json:"sessionId"`
	PhoneNumber string  `json:"phoneNumber"`
	Rating      int     `json:"rating"`
	Comment     *string `json:"comment,omitempty"`
	CreatedAt   string  `json:"createdAt"`
}

func ToFeedbackResponse(feedback *entity.Feedback) FeedbackResponse {
	return FeedbackResponse{
		ID:          feedback.ID.String(),
		SessionID:   feedback.SessionID.String(),
		PhoneNumber: feedback.PhoneNumber,
		Rating:      feedback.Rating,
		Comment:     feedback.Comment,
		CreatedAt:   feedback.CreatedAt.Format(time.RFC3339),
	}
}

type CreateFeedbackRequest struct {
	SessionID   string  `json:"sessionId" validate:"required,uuid"`
	PhoneNumber string  `json:"phoneNumber" validate:"required,min=10,max=20"`
	Rating      int     `json:"rating" validate:"required,min=1,max=5"`
	Comment     *string `json:"comment,omitempty" validate:"omitempty,max=1000"`
}

type CreateFeedbackResponse struct {
	ID string `json:"id"`
}

type GetFeedbacksQuery struct {
	Page        int     `query:"page" validate:"omitempty,min=1"`
	Limit       int     `query:"limit" validate:"omitempty,min=1,max=100"`
	PhoneNumber *string `query:"phoneNumber" validate:"omitempty,max=20"`
	MinRating   *int    `query:"minRating" validate:"omitempty,min=1,max=5"`
	MaxRating   *int    `query:"maxRating" validate:"omitempty,min=1,max=5"`
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
