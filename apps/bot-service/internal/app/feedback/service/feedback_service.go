package service

import (
	"context"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	"github.com/google/uuid"
)

func (s *FeedbackService) Create(ctx context.Context, req *dto.CreateFeedbackRequest) (*dto.CreateFeedbackResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	userID, err := s.uuidPkg.Parse(req.UserID)
	if err != nil {
		return nil, errx.ErrUserNotFound.WithDetails(map[string]any{
			"user_id": req.UserID,
		}).WithLocation("FeedbackService.Create").WithError(err)
	}

	id, err := s.uuidPkg.NewV7()
	if err != nil {
		return nil, errx.ErrInternalServer.WithLocation("FeedbackService.Create").WithError(err)
	}

	feedback := &entity.Feedback{
		ID:        id,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
	}

	if err := s.feedbackRepo.Create(ctx, feedback); err != nil {
		return nil, err
	}

	res := dto.CreateFeedbackResponse{
		ID: id.String(),
	}

	return &res, nil
}

func (s *FeedbackService) GetByID(ctx context.Context, param *dto.GetFeedbackByIDParam) (*dto.GetFeedbackByIDResponse, error) {
	if err := s.validator.Validate(param); err != nil {
		return nil, err
	}

	id, err := s.uuidPkg.Parse(param.ID)
	if err != nil {
		return nil, errx.ErrFeedbackNotFound.WithDetails(map[string]any{
			"id": param.ID,
		}).WithLocation("FeedbackService.GetByID").WithError(err)
	}

	feedback, err := s.feedbackRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	feedbackRes := dto.ToFeedbackResponse(feedback)

	res := &dto.GetFeedbackByIDResponse{
		Feedback: feedbackRes,
	}

	return res, nil
}

func (s *FeedbackService) List(ctx context.Context, query *dto.GetFeedbacksQuery) (*dto.GetFeedbacksResponse, error) {
	if err := s.validator.Validate(query); err != nil {
		return nil, err
	}

	limit := min(max(query.Limit, 10), 100)
	page := max(query.Page, 1)

	var userID *uuid.UUID
	if query.UserID != nil {
		parsedUserID, err := s.uuidPkg.Parse(*query.UserID)
		if err != nil {
			return nil, errx.ErrUserNotFound.WithDetails(map[string]any{
				"user_id": *query.UserID,
			}).WithLocation("FeedbackService.List").WithError(err)
		}
		userID = &parsedUserID
	}

	filter := entity.GetFeedbacksFilter{
		Offset:    (page - 1) * limit,
		Limit:     limit,
		UserID:    userID,
		MinRating: query.MinRating,
		MaxRating: query.MaxRating,
	}

	feedbacks, total, err := s.feedbackRepo.List(ctx, &filter)
	if err != nil {
		return nil, err
	}

	feedbackResponses := make([]dto.FeedbackResponse, 0, len(feedbacks))
	for i := range feedbacks {
		feedbackResponses = append(feedbackResponses, dto.ToFeedbackResponse(&feedbacks[i]))
	}

	paginationResponse := dto.NewPaginationResponse(total, page, limit)

	res := &dto.GetFeedbacksResponse{
		Feedbacks: feedbackResponses,
	}

	res.Meta.Pagination = paginationResponse

	return res, nil
}

func (s *FeedbackService) GetMetrics(ctx context.Context) (*dto.GetFeedbackMetricsResponse, error) {
	satisfactionScore, err := s.feedbackRepo.GetMetrics(ctx)
	if err != nil {
		return nil, err
	}

	res := &dto.GetFeedbackMetricsResponse{
		SatisfactionScore: satisfactionScore,
	}

	return res, nil
}

func (s *FeedbackService) GetSatisfactionTrend(ctx context.Context, query *dto.GetSatisfactionTrendQuery) (*dto.GetSatisfactionTrendResponse, error) {
	if err := s.validator.Validate(query); err != nil {
		return nil, err
	}

	days := 30
	if query.Days > 0 {
		days = query.Days
	}

	results, err := s.feedbackRepo.GetSatisfactionTrend(ctx, days)
	if err != nil {
		return nil, err
	}

	trend := make([]dto.SatisfactionTrendData, 0, len(results))
	for _, result := range results {
		date, ok := result["date"].(time.Time)
		if !ok {
			continue
		}

		avgSatisfaction := 0.0
		if val, ok := result["avg_satisfaction"].(float64); ok {
			avgSatisfaction = val
		}

		trend = append(trend, dto.SatisfactionTrendData{
			Date:            date.Format(time.RFC3339),
			AvgSatisfaction: avgSatisfaction,
		})
	}

	res := &dto.GetSatisfactionTrendResponse{
		Trend: trend,
	}

	return res, nil
}
