package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	userRepoMock "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/user/repository/mock"
	mockUUID "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid/mock"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
	mockValidator "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := userRepoMock.NewMockUserRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewUserService(mockUserRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testID := uuid.New()
	assignedTo := "John Doe"
	notes := "Test notes"

	tests := []struct {
		name    string
		req     *dto.CreateUserRequest
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name: "success",
			req: &dto.CreateUserRequest{
				PhoneNumber: "+1234567890",
				Label:       "Test User",
				AssignedTo:  &assignedTo,
				Notes:       &notes,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().NewV7().Return(testID, nil)
				mockUserRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entity.User) error {
					assert.Equal(t, testID, user.ID)
					assert.Equal(t, "+1234567890", user.PhoneNumber)
					assert.Equal(t, "Test User", user.Label)
					assert.Equal(t, &assignedTo, user.AssignedTo)
					assert.Equal(t, &notes, user.Notes)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "validation error",
			req: &dto.CreateUserRequest{
				PhoneNumber: "123",
				Label:       "Test User",
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
			name: "uuid generation error",
			req: &dto.CreateUserRequest{
				PhoneNumber: "+1234567890",
				Label:       "Test User",
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().NewV7().Return(uuid.Nil, errors.New("uuid generation failed"))
			},
			wantErr: true,
		},
		{
			name: "duplicate phone number",
			req: &dto.CreateUserRequest{
				PhoneNumber: "+1234567890",
				Label:       "Test User",
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().NewV7().Return(testID, nil)
				mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(errx.ErrUserPhoneExists)
			},
			wantErr: true,
			errType: errx.ErrUserPhoneExists,
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

func TestUserService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := userRepoMock.NewMockUserRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewUserService(mockUserRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testID := uuid.New()
	testUser := &entity.User{
		ID:          testID,
		PhoneNumber: "+1234567890",
		Label:       "Test User",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		param   *dto.GetUserByIDParam
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name: "success",
			param: &dto.GetUserByIDParam{
				ID: testID.String(),
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().FindByID(ctx, testID).Return(testUser, nil)
			},
			wantErr: false,
		},
		{
			name: "validation error",
			param: &dto.GetUserByIDParam{
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
			param: &dto.GetUserByIDParam{
				ID: "invalid-uuid",
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse("invalid-uuid").Return(uuid.Nil, errors.New("invalid uuid"))
			},
			wantErr: true,
		},
		{
			name: "user not found",
			param: &dto.GetUserByIDParam{
				ID: testID.String(),
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().FindByID(ctx, testID).Return(nil, errx.ErrUserNotFound)
			},
			wantErr: true,
			errType: errx.ErrUserNotFound,
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
				assert.Equal(t, testID.String(), result.User.ID)
				assert.Equal(t, "+1234567890", result.User.PhoneNumber)
				assert.Equal(t, "Test User", result.User.Label)
			}
		})
	}
}

func TestUserService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := userRepoMock.NewMockUserRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewUserService(mockUserRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testUsers := []entity.User{
		{
			ID:          uuid.New(),
			PhoneNumber: "+1234567890",
			Label:       "User 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			PhoneNumber: "+0987654321",
			Label:       "User 2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	assignedTo := "John Doe"

	tests := []struct {
		name      string
		query     *dto.GetUsersQuery
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
			query: &dto.GetUsersQuery{
				Page:  0,
				Limit: 0,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetUsersFilter) ([]entity.User, int64, error) {
					assert.Equal(t, 0, filter.Offset)
					assert.Equal(t, 10, filter.Limit)
					return testUsers, 2, nil
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
			query: &dto.GetUsersQuery{
				Page:  2,
				Limit: 20,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetUsersFilter) ([]entity.User, int64, error) {
					assert.Equal(t, 20, filter.Offset)
					assert.Equal(t, 20, filter.Limit)
					return testUsers, int64(40), nil
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
			name: "success with search filter",
			query: &dto.GetUsersQuery{
				Page:   1,
				Limit:  10,
				Search: "test",
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetUsersFilter) ([]entity.User, int64, error) {
					assert.Equal(t, "test", filter.Search)
					return testUsers, int64(2), nil
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
			name: "success with assignedTo filter",
			query: &dto.GetUsersQuery{
				Page:       1,
				Limit:      10,
				AssignedTo: &assignedTo,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().List(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, filter *entity.GetUsersFilter) ([]entity.User, int64, error) {
					assert.Equal(t, &assignedTo, filter.AssignedTo)
					return testUsers, int64(2), nil
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
			query: &dto.GetUsersQuery{
				Page:  1,
				Limit: 10,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().List(ctx, gomock.Any()).Return([]entity.User{}, int64(0), nil)
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
			query: &dto.GetUsersQuery{
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
				assert.Len(t, result.Users, tt.wantCount)
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

func TestUserService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := userRepoMock.NewMockUserRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewUserService(mockUserRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testID := uuid.New()
	testUser := &entity.User{
		ID:          testID,
		PhoneNumber: "+1234567890",
		Label:       "Test User",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	newPhone := "+0987654321"
	newLabel := "Updated User"
	newAssignedTo := "Jane Doe"
	newNotes := "Updated notes"

	tests := []struct {
		name    string
		param   *dto.UpdateUserParam
		req     *dto.UpdateUserRequest
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name: "success - update phone number",
			param: &dto.UpdateUserParam{
				ID: testID.String(),
			},
			req: &dto.UpdateUserRequest{
				PhoneNumber: &newPhone,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(2)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().FindByID(ctx, testID).Return(testUser, nil)
				mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entity.User) error {
					assert.Equal(t, newPhone, user.PhoneNumber)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "success - update label",
			param: &dto.UpdateUserParam{
				ID: testID.String(),
			},
			req: &dto.UpdateUserRequest{
				Label: &newLabel,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(2)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().FindByID(ctx, testID).Return(testUser, nil)
				mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entity.User) error {
					assert.Equal(t, newLabel, user.Label)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "success - update multiple fields",
			param: &dto.UpdateUserParam{
				ID: testID.String(),
			},
			req: &dto.UpdateUserRequest{
				PhoneNumber: &newPhone,
				Label:       &newLabel,
				AssignedTo:  &newAssignedTo,
				Notes:       &newNotes,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(2)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().FindByID(ctx, testID).Return(testUser, nil)
				mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entity.User) error {
					assert.Equal(t, newPhone, user.PhoneNumber)
					assert.Equal(t, newLabel, user.Label)
					assert.Equal(t, &newAssignedTo, user.AssignedTo)
					assert.Equal(t, &newNotes, user.Notes)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "param validation error",
			param: &dto.UpdateUserParam{
				ID: "invalid",
			},
			req: &dto.UpdateUserRequest{
				Label: &newLabel,
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
			name: "request validation error",
			param: &dto.UpdateUserParam{
				ID: testID.String(),
			},
			req: &dto.UpdateUserRequest{
				PhoneNumber: &newPhone,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(&dto.UpdateUserParam{ID: testID.String()}).Return(nil)
				mockValidator.EXPECT().Validate(gomock.Any()).Return(validator.ValidationErrors{
					"body": validator.ValidationError{
						Message: "validation error",
					},
				})
			},
			wantErr: true,
		},
		{
			name: "invalid uuid",
			param: &dto.UpdateUserParam{
				ID: "invalid-uuid",
			},
			req: &dto.UpdateUserRequest{
				Label: &newLabel,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(2)
				mockUUID.EXPECT().Parse("invalid-uuid").Return(uuid.Nil, errors.New("invalid uuid"))
			},
			wantErr: true,
		},
		{
			name: "user not found",
			param: &dto.UpdateUserParam{
				ID: testID.String(),
			},
			req: &dto.UpdateUserRequest{
				Label: &newLabel,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(2)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().FindByID(ctx, testID).Return(nil, errx.ErrUserNotFound)
			},
			wantErr: true,
			errType: errx.ErrUserNotFound,
		},
		{
			name: "duplicate phone number",
			param: &dto.UpdateUserParam{
				ID: testID.String(),
			},
			req: &dto.UpdateUserRequest{
				PhoneNumber: &newPhone,
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(2)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().FindByID(ctx, testID).Return(testUser, nil)
				mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(errx.ErrUserPhoneExists)
			},
			wantErr: true,
			errType: errx.ErrUserPhoneExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.Update(ctx, tt.param, tt.req)

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

func TestUserService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := userRepoMock.NewMockUserRepository(ctrl)
	mockValidator := mockValidator.NewMockCustomValidatorInterface(ctrl)
	mockUUID := mockUUID.NewMockUUIDInterface(ctrl)

	service := NewUserService(mockUserRepo, mockValidator, mockUUID)
	ctx := context.Background()

	testID := uuid.New()

	tests := []struct {
		name    string
		param   *dto.DeleteUserParam
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name: "success",
			param: &dto.DeleteUserParam{
				ID: testID.String(),
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().Delete(ctx, testID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error",
			param: &dto.DeleteUserParam{
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
			param: &dto.DeleteUserParam{
				ID: "invalid-uuid",
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse("invalid-uuid").Return(uuid.Nil, errors.New("invalid uuid"))
			},
			wantErr: true,
		},
		{
			name: "user not found",
			param: &dto.DeleteUserParam{
				ID: testID.String(),
			},
			setup: func() {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockUUID.EXPECT().Parse(testID.String()).Return(testID, nil)
				mockUserRepo.EXPECT().Delete(ctx, testID).Return(errx.ErrUserNotFound)
			},
			wantErr: true,
			errType: errx.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.Delete(ctx, tt.param)

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
