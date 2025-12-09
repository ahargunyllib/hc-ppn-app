package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	feedbackRepoMock "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/repository/mock"
	mockUUID "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid/mock"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
	mockValidator "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFeedbackService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedbackRepo := feedbackRepoMock.NewMockFeedbackRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewFeedbackService(mockFeedbackRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testID := uuid.New()
	testUserID := uuid.New()
	comment := "Great service!"

	tests := []struct {
		name    string
		req     *dto.CreateFeedbackRequest
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name: "success with comment",
			req: &dto.CreateFeedbackRequest{
				UserID:  testUserID.String(),
				Rating:  5,
				Comment: &comment,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testUserID.String()).Return(testUserID, nil)
				mockUUID.EXPECT().NewV7().Return(testID, nil)
				mockFeedbackRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, feedback *entity.Feedback) error {
					assert.Equal(t, testID, feedback.ID)
					assert.Equal(t, testUserID, feedback.UserID)
					assert.Equal(t, 5, feedback.Rating)
					assert.Equal(t, &comment, feedback.Comment)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "success without comment",
			req: &dto.CreateFeedbackRequest{
				UserID: testUserID.String(),
				Rating: 4,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testUserID.String()).Return(testUserID, nil)
				mockUUID.EXPECT().NewV7().Return(testID, nil)
				mockFeedbackRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, feedback *entity.Feedback) error {
					assert.Equal(t, testID, feedback.ID)
					assert.Equal(t, testUserID, feedback.UserID)
					assert.Equal(t, 4, feedback.Rating)
					assert.Nil(t, feedback.Comment)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "validation error",
			req: &dto.CreateFeedbackRequest{
				UserID: "invalid",
				Rating: 6,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(validator.ValidationErrors{
					"body": validator.ValidationError{
						Message: "validation error",
					},
				})
			},
			wantErr: true,
		},
		{
			name: "invalid user id",
			req: &dto.CreateFeedbackRequest{
				UserID: "invalid-uuid",
				Rating: 5,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse("invalid-uuid").Return(uuid.Nil, errors.New("invalid uuid"))
			},
			wantErr: true,
		},
		{
			name: "uuid generation error",
			req: &dto.CreateFeedbackRequest{
				UserID: testUserID.String(),
				Rating: 5,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testUserID.String()).Return(testUserID, nil)
				mockUUID.EXPECT().NewV7().Return(uuid.Nil, errors.New("uuid generation failed"))
			},
			wantErr: true,
		},
		{
			name: "user not found",
			req: &dto.CreateFeedbackRequest{
				UserID: testUserID.String(),
				Rating: 5,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testUserID.String()).Return(testUserID, nil)
				mockUUID.EXPECT().NewV7().Return(testID, nil)
				mockFeedbackRepo.EXPECT().Create(ctx, gomock.Any()).Return(errx.ErrUserNotFound)
			},
			wantErr: true,
			errType: errx.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := service.Create(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testID.String(), result.ID)
			}
		})
	}
}

func TestFeedbackService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedbackRepo := feedbackRepoMock.NewMockFeedbackRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewFeedbackService(mockFeedbackRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testID := uuid.New()
	testUserID := uuid.New()
	comment := "Great service!"
	testFeedback := &entity.Feedback{
		ID:        testID,
		UserID:    testUserID,
		Rating:    5,
		Comment:   &comment,
		CreatedAt: time.Now(),
		User: entity.User{
			ID:          testUserID,
			PhoneNumber: "+1234567890",
			Name:        "Test User",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	tests := []struct {
		name    string
		param   *dto.GetFeedbackByIDParam
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name: "success",
			param: &dto.GetFeedbackByIDParam{
				ID: testID.String(),
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockFeedbackRepo.EXPECT().FindByID(ctx, testID).Return(testFeedback, nil)
			},
			wantErr: false,
		},
		{
			name: "validation error",
			param: &dto.GetFeedbackByIDParam{
				ID: "invalid",
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(validator.ValidationErrors{
					"param": validator.ValidationError{
						Message: "validation error",
					},
				})
			},
			wantErr: true,
		},
		{
			name: "invalid uuid",
			param: &dto.GetFeedbackByIDParam{
				ID: "invalid-uuid",
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse("invalid-uuid").Return(uuid.Nil, errors.New("invalid uuid"))
			},
			wantErr: true,
		},
		{
			name: "feedback not found",
			param: &dto.GetFeedbackByIDParam{
				ID: testID.String(),
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockFeedbackRepo.EXPECT().FindByID(ctx, testID).Return(nil, errx.ErrFeedbackNotFound)
			},
			wantErr: true,
			errType: errx.ErrFeedbackNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := service.GetByID(ctx, tt.param)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testID.String(), result.Feedback.ID)
				assert.Equal(t, 5, result.Feedback.Rating)
				assert.Equal(t, &comment, result.Feedback.Comment)
			}
		})
	}
}

func TestFeedbackService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedbackRepo := feedbackRepoMock.NewMockFeedbackRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewFeedbackService(mockFeedbackRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testUserID := uuid.New()
	testFeedbacks := []entity.Feedback{
		{
			ID:        uuid.New(),
			UserID:    testUserID,
			Rating:    5,
			CreatedAt: time.Now(),
			User: entity.User{
				ID:          testUserID,
				PhoneNumber: "+1234567890",
				Name:        "User 1",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Rating:    4,
			CreatedAt: time.Now(),
			User: entity.User{
				ID:          uuid.New(),
				PhoneNumber: "+0987654321",
				Name:        "User 2",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	userIDStr := testUserID.String()
	minRating := 3
	maxRating := 5

	tests := []struct {
		name      string
		query     *dto.GetFeedbacksQuery
		setup     func()
		wantErr   bool
		wantCount int
		wantPage  int
		wantLimit int
		wantTotal int64
		wantPages int
	}{
		{
			name: "success with default pagination",
			query: &dto.GetFeedbacksQuery{
				Page:  0,
				Limit: 0,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockFeedbackRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetFeedbacksFilter) ([]entity.Feedback, int64, error) {
					assert.Equal(t, 0, filter.Offset)
					assert.Equal(t, 10, filter.Limit)
					return testFeedbacks, 2, nil
				})
			},
			wantErr:   false,
			wantCount: 2,
			wantPage:  1,
			wantLimit: 10,
			wantTotal: 2,
			wantPages: 1,
		},
		{
			name: "success with custom pagination",
			query: &dto.GetFeedbacksQuery{
				Page:  2,
				Limit: 20,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockFeedbackRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetFeedbacksFilter) ([]entity.Feedback, int64, error) {
					assert.Equal(t, 20, filter.Offset)
					assert.Equal(t, 20, filter.Limit)
					return testFeedbacks, int64(40), nil
				})
			},
			wantErr:   false,
			wantCount: 2,
			wantPage:  2,
			wantLimit: 20,
			wantTotal: 40,
			wantPages: 2,
		},
		{
			name: "success with userID filter",
			query: &dto.GetFeedbacksQuery{
				Page:   1,
				Limit:  10,
				UserID: &userIDStr,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(userIDStr).Return(testUserID, nil)
				mockFeedbackRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetFeedbacksFilter) ([]entity.Feedback, int64, error) {
					assert.NotNil(t, filter.UserID)
					assert.Equal(t, testUserID, *filter.UserID)
					return testFeedbacks, int64(2), nil
				})
			},
			wantErr:   false,
			wantCount: 2,
			wantPage:  1,
			wantLimit: 10,
			wantTotal: 2,
			wantPages: 1,
		},
		{
			name: "success with minRating filter",
			query: &dto.GetFeedbacksQuery{
				Page:      1,
				Limit:     10,
				MinRating: &minRating,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockFeedbackRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetFeedbacksFilter) ([]entity.Feedback, int64, error) {
					assert.NotNil(t, filter.MinRating)
					assert.Equal(t, minRating, *filter.MinRating)
					return testFeedbacks, int64(2), nil
				})
			},
			wantErr:   false,
			wantCount: 2,
			wantPage:  1,
			wantLimit: 10,
			wantTotal: 2,
			wantPages: 1,
		},
		{
			name: "success with maxRating filter",
			query: &dto.GetFeedbacksQuery{
				Page:      1,
				Limit:     10,
				MaxRating: &maxRating,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockFeedbackRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetFeedbacksFilter) ([]entity.Feedback, int64, error) {
					assert.NotNil(t, filter.MaxRating)
					assert.Equal(t, maxRating, *filter.MaxRating)
					return testFeedbacks, int64(2), nil
				})
			},
			wantErr:   false,
			wantCount: 2,
			wantPage:  1,
			wantLimit: 10,
			wantTotal: 2,
			wantPages: 1,
		},
		{
			name: "success with multiple filters",
			query: &dto.GetFeedbacksQuery{
				Page:      1,
				Limit:     10,
				UserID:    &userIDStr,
				MinRating: &minRating,
				MaxRating: &maxRating,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(userIDStr).Return(testUserID, nil)
				mockFeedbackRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetFeedbacksFilter) ([]entity.Feedback, int64, error) {
					assert.NotNil(t, filter.UserID)
					assert.Equal(t, testUserID, *filter.UserID)
					assert.NotNil(t, filter.MinRating)
					assert.Equal(t, minRating, *filter.MinRating)
					assert.NotNil(t, filter.MaxRating)
					assert.Equal(t, maxRating, *filter.MaxRating)
					return testFeedbacks, int64(2), nil
				})
			},
			wantErr:   false,
			wantCount: 2,
			wantPage:  1,
			wantLimit: 10,
			wantTotal: 2,
			wantPages: 1,
		},
		{
			name: "empty results",
			query: &dto.GetFeedbacksQuery{
				Page:  1,
				Limit: 10,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockFeedbackRepo.EXPECT().List(ctx, gomock.Any()).Return([]entity.Feedback{}, int64(0), nil)
			},
			wantErr:   false,
			wantCount: 0,
			wantPage:  1,
			wantLimit: 10,
			wantTotal: 0,
			wantPages: 0,
		},
		{
			name: "validation error",
			query: &dto.GetFeedbacksQuery{
				Page:  -1,
				Limit: 200,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(validator.ValidationErrors{
					"query": validator.ValidationError{
						Message: "validation error",
					},
				})
			},
			wantErr: true,
		},
		{
			name: "invalid user id in filter",
			query: &dto.GetFeedbacksQuery{
				Page:   1,
				Limit:  10,
				UserID: &userIDStr,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(userIDStr).Return(uuid.Nil, errors.New("invalid uuid"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := service.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Feedbacks, tt.wantCount)
				if tt.wantPage > 0 {
					assert.Equal(t, tt.wantPage, result.Meta.Pagination.Page)
				}
				if tt.wantLimit > 0 {
					assert.Equal(t, tt.wantLimit, result.Meta.Pagination.Limit)
				}
				if tt.wantTotal >= 0 {
					assert.Equal(t, tt.wantTotal, result.Meta.Pagination.TotalData)
				}
				if tt.wantPages >= 0 {
					assert.Equal(t, tt.wantPages, result.Meta.Pagination.TotalPage)
				}
			}
		})
	}
}

func TestFeedbackService_GetMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedbackRepo := feedbackRepoMock.NewMockFeedbackRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewFeedbackService(mockFeedbackRepo, mockValidator, mockUUID)
	ctx := context.Background()

	tests := []struct {
		name             string
		setup            func()
		wantErr          bool
		wantSatisfaction float64
		wantTotal        int
	}{
		{
			name: "success with satisfaction score and total feedbacks",
			setup: func() {
				mockFeedbackRepo.EXPECT().GetMetrics(ctx).Return(85.5, 42, nil)
			},
			wantErr:          false,
			wantSatisfaction: 85.5,
			wantTotal:        42,
		},
		{
			name: "success with zero satisfaction score and zero feedbacks",
			setup: func() {
				mockFeedbackRepo.EXPECT().GetMetrics(ctx).Return(0.0, 0, nil)
			},
			wantErr:          false,
			wantSatisfaction: 0.0,
			wantTotal:        0,
		},
		{
			name: "success with high feedback count",
			setup: func() {
				mockFeedbackRepo.EXPECT().GetMetrics(ctx).Return(92.3, 1500, nil)
			},
			wantErr:          false,
			wantSatisfaction: 92.3,
			wantTotal:        1500,
		},
		{
			name: "repository error",
			setup: func() {
				mockFeedbackRepo.EXPECT().GetMetrics(ctx).Return(0.0, 0, errx.ErrInternalServer)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := service.GetMetrics(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantSatisfaction, result.SatisfactionScore)
				assert.Equal(t, tt.wantTotal, result.TotalFeedbacks)
			}
		})
	}
}

func TestFeedbackService_GetSatisfactionTrend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedbackRepo := feedbackRepoMock.NewMockFeedbackRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewFeedbackService(mockFeedbackRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testDate := time.Date(2025, 12, 6, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		setup      func()
		wantErr    bool
		wantCount  int
		checkFirst func(*testing.T, *dto.GetSatisfactionTrendResponse)
	}{
		{
			name: "success with trend data",
			setup: func() {
				mockFeedbackRepo.EXPECT().GetSatisfactionTrend(ctx).Return([]entity.SatisfactionTrendRow{
					{Date: testDate, AvgSatisfaction: 4.5},
					{Date: testDate.AddDate(0, 0, 1), AvgSatisfaction: 4.2},
					{Date: testDate.AddDate(0, 0, 2), AvgSatisfaction: 4.8},
				}, nil)
			},
			wantErr:   false,
			wantCount: 3,
			checkFirst: func(t *testing.T, res *dto.GetSatisfactionTrendResponse) {
				assert.Equal(t, testDate.Format(time.RFC3339), res.Trend[0].Date)
				assert.Equal(t, 4.5, res.Trend[0].AvgSatisfaction)
			},
		},
		{
			name: "success with empty trend data",
			setup: func() {
				mockFeedbackRepo.EXPECT().GetSatisfactionTrend(ctx).Return([]entity.SatisfactionTrendRow{}, nil)
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name: "repository error",
			setup: func() {
				mockFeedbackRepo.EXPECT().GetSatisfactionTrend(ctx).Return(nil, errx.ErrInternalServer)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := service.GetSatisfactionTrend(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Trend, tt.wantCount)
				if tt.checkFirst != nil && tt.wantCount > 0 {
					tt.checkFirst(t, result)
				}
			}
		})
	}
}
