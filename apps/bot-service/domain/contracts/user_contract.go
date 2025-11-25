package contracts

import (
	"context"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=../../internal/app/user/repository/mock/mock_user_repository.go -package=mock github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/contracts UserRepository

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	List(ctx context.Context, filter *entity.GetUsersFilter) ([]entity.User, int64, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserService interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	GetByID(ctx context.Context, param *dto.GetUserByIDParam) (*dto.GetUserByIDResponse, error)
	List(ctx context.Context, query *dto.GetUsersQuery) (*dto.GetUsersResponse, error)
	Update(ctx context.Context, param *dto.UpdateUserParam, req *dto.UpdateUserRequest) error
	Delete(ctx context.Context, param *dto.DeleteUserParam) error
}
