package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	topicRepoMock "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/topic/repository/mock"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
	mockValidator "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTopicService_BulkCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopicRepo := topicRepoMock.NewMockTopicRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)

	service := NewTopicService(mockTopicRepo, mockValidator)
	ctx := context.Background()

	tests := []struct {
		name         string
		req          *dto.BulkCreateTopicsRequest
		setup        func()
		wantErr      bool
		errType      error
	}{
		{
			name: "success - create topics",
			req: &dto.BulkCreateTopicsRequest{
				Topics: []dto.TopicRequest{
					{Title: "Billing", Count: 5},
					{Title: "Support", Count: 3},
				},
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockTopicRepo.EXPECT().BulkCreate(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, topics []entity.Topic) error {
					assert.Len(t, topics, 2)
					assert.Equal(t, "Billing", topics[0].Title)
					assert.Equal(t, 5, topics[0].Count)
					assert.Equal(t, "Support", topics[1].Title)
					assert.Equal(t, 3, topics[1].Count)
					return nil
				})
			},
			wantErr:      false,
		},
		{
			name: "success - single topic",
			req: &dto.BulkCreateTopicsRequest{
				Topics: []dto.TopicRequest{
					{Title: "General Inquiry", Count: 10},
				},
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockTopicRepo.EXPECT().BulkCreate(ctx, gomock.Any()).Return(nil)
			},
			wantErr:      false,
		},
		{
			name: "success - allow duplicates",
			req: &dto.BulkCreateTopicsRequest{
				Topics: []dto.TopicRequest{
					{Title: "Billing", Count: 5},
					{Title: "Billing", Count: 3},
					{Title: "Support", Count: 2},
				},
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockTopicRepo.EXPECT().BulkCreate(ctx, gomock.Any()).Return(nil)
			},
			wantErr:      false,
		},
		{
			name: "validation error - empty topics",
			req: &dto.BulkCreateTopicsRequest{
				Topics: []dto.TopicRequest{},
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(validator.ValidationErrors{
					"body.topics": validator.ValidationError{
						Message: "topics is required and must have at least 1 item",
					},
				})
			},
			wantErr: true,
		},
		{
			name: "validation error - invalid title",
			req: &dto.BulkCreateTopicsRequest{
				Topics: []dto.TopicRequest{
					{Title: "", Count: 5},
				},
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(validator.ValidationErrors{
					"body.topics[0].title": validator.ValidationError{
						Message: "title is required",
					},
				})
			},
			wantErr: true,
		},
		{
			name: "validation error - invalid count",
			req: &dto.BulkCreateTopicsRequest{
				Topics: []dto.TopicRequest{
					{Title: "Billing", Count: 0},
				},
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(validator.ValidationErrors{
					"body.topics[0].count": validator.ValidationError{
						Message: "count must be at least 1",
					},
				})
			},
			wantErr: true,
		},
		{
			name: "validation error - too many topics",
			req: &dto.BulkCreateTopicsRequest{
				Topics: make([]dto.TopicRequest, 101),
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(validator.ValidationErrors{
					"body.topics": validator.ValidationError{
						Message: "topics must have at most 100 items",
					},
				})
			},
			wantErr: true,
		},
		{
			name: "repository error",
			req: &dto.BulkCreateTopicsRequest{
				Topics: []dto.TopicRequest{
					{Title: "Billing", Count: 5},
				},
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockTopicRepo.EXPECT().BulkCreate(ctx, gomock.Any()).Return(errx.ErrInternalServer)
			},
			wantErr: true,
			errType: errx.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.BulkCreate(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTopicService_GetHotTopics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTopicRepo := topicRepoMock.NewMockTopicRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)

	service := NewTopicService(mockTopicRepo, mockValidator)
	ctx := context.Background()

	testDate1 := time.Now().AddDate(0, 0, -2)
	testDate2 := time.Now().AddDate(0, 0, -1)
	testDate3 := time.Now()

	testHotTopics := []entity.Topic{
		{ID: 1, Title: "Billing", Count: 18, CreatedAt: testDate1},
		{ID: 2, Title: "Support", Count: 7, CreatedAt: testDate2},
		{ID: 3, Title: "Feedback", Count: 5, CreatedAt: testDate3},
		{ID: 4, Title: "Technical Issue", Count: 3, CreatedAt: testDate1},
		{ID: 5, Title: "General Inquiry", Count: 2, CreatedAt: testDate2},
	}

	tests := []struct {
		name      string
		setup     func()
		wantErr   bool
		wantCount int
		errType   error
	}{
		{
			name: "success with hot topics data",
			setup: func() {
				mockTopicRepo.EXPECT().GetHotTopics(ctx).Return(testHotTopics, nil)
			},
			wantErr:   false,
			wantCount: 5,
		},
		{
			name: "success with empty results",
			setup: func() {
				mockTopicRepo.EXPECT().GetHotTopics(ctx).Return([]entity.Topic{}, nil)
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name: "repository error",
			setup: func() {
				mockTopicRepo.EXPECT().GetHotTopics(ctx).Return(nil, errx.ErrInternalServer)
			},
			wantErr: true,
			errType: errx.ErrInternalServer,
		},
		{
			name: "database connection error",
			setup: func() {
				mockTopicRepo.EXPECT().GetHotTopics(ctx).Return(nil, errors.New("database connection failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := service.GetHotTopics(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Topics, tt.wantCount)

				if tt.wantCount > 0 {
					for _, topic := range result.Topics {
						assert.NotZero(t, topic.ID)
						assert.NotEmpty(t, topic.Title)
						assert.GreaterOrEqual(t, topic.Count, 0)
						assert.NotEmpty(t, topic.CreatedAt)
					}
				}
			}
		})
	}
}
