package service

import (
	"context"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
)

func (s *FeedbackService) Create(ctx context.Context, req *dto.CreateFeedbackRequest) (*dto.CreateFeedbackResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	sessionID, err := s.uuidPkg.Parse(req.SessionID)
	if err != nil {
		return nil, errx.ErrSessionNotFound.WithDetails(map[string]any{
			"session_id": req.SessionID,
		}).WithLocation("FeedbackService.Create").WithError(err)
	}

	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.Status == entity.SessionStatusClosed {
		return nil, errx.ErrSessionAlreadyClosed.WithLocation("FeedbackService.Create")
	}

	existingFeedback, err := s.feedbackRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if existingFeedback != nil {
		return nil, errx.ErrFeedbackAlreadyExists.WithLocation("FeedbackService.Create")
	}

	id, err := s.uuidPkg.NewV7()
	if err != nil {
		return nil, errx.ErrInternalServer.WithLocation("FeedbackService.Create").WithError(err)
	}

	feedback := &entity.Feedback{
		ID:          id,
		SessionID:   sessionID,
		PhoneNumber: req.PhoneNumber,
		Rating:      req.Rating,
		Comment:     req.Comment,
		CreatedAt:   time.Now(),
	}

	if err := s.feedbackRepo.Create(ctx, feedback); err != nil {
		return nil, err
	}

	session.Status = entity.SessionStatusClosed
	session.UpdatedAt = time.Now()
	if err := s.sessionRepo.Update(ctx, session); err != nil {
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

	filter := entity.GetFeedbacksFilter{
		Offset:      (page - 1) * limit,
		Limit:       limit,
		PhoneNumber: query.PhoneNumber,
		MinRating:   query.MinRating,
		MaxRating:   query.MaxRating,
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
