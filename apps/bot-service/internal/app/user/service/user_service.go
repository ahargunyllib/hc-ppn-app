package service

import (
	"context"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
)

func (s *UserService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	id, err := s.uuidPkg.NewV7()
	if err != nil {
		return nil, errx.ErrInternalServer.WithLocation("UserService.Create").WithError(err)
	}

	user := &entity.User{
		ID:          id,
		PhoneNumber: req.PhoneNumber,
		Label:       req.Label,
		AssignedTo:  req.AssignedTo,
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	res := dto.CreateUserResponse{
		ID: id.String(),
	}

	return &res, nil
}

func (s *UserService) GetByID(ctx context.Context, param *dto.GetUserByIDParam) (*dto.GetUserByIDResponse, error) {
	// Validate param
	if err := s.validator.Validate(param); err != nil {
		return nil, err
	}

	id, err := s.uuidPkg.Parse(param.ID)
	if err != nil {
		return nil, errx.ErrUserNotFound.WithDetails(map[string]any{
			"id": param.ID,
		}).WithLocation("UserService.GetByID").WithError(err)
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userRes := dto.ToUserResponse(user)

	res := &dto.GetUserByIDResponse{
		User: userRes,
	}

	return res, nil
}

func (s *UserService) List(ctx context.Context, query *dto.GetUsersQuery) (*dto.GetUsersResponse, error) {
	if err := s.validator.Validate(query); err != nil {
		return nil, err
	}

	limit := min(max(query.Limit, 10), 100)
	page := max(query.Page, 1)

	filter := entity.GetUsersFilter{
		Offset:     (page - 1) * limit,
		Limit:      limit,
		Search:     query.Search,
		AssignedTo: query.AssignedTo,
	}

	users, total, err := s.userRepo.List(ctx, &filter)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		userResponses = append(userResponses, dto.ToUserResponse(&users[i]))
	}

	paginationResponse := dto.NewPaginationResponse(total, page, limit)

	res := &dto.GetUsersResponse{
		Users: userResponses,
	}

	res.Meta.Pagination = paginationResponse

	return res, nil
}

func (s *UserService) Update(ctx context.Context, param *dto.UpdateUserParam, req *dto.UpdateUserRequest) error {
	if err := s.validator.Validate(param); err != nil {
		return err
	}

	if err := s.validator.Validate(req); err != nil {
		return err
	}

	id, err := s.uuidPkg.Parse(param.ID)
	if err != nil {
		return errx.ErrUserNotFound.WithDetails(map[string]any{
			"id": param.ID,
		}).WithLocation("UserService.Update").WithError(err)
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}
	if req.Label != nil {
		user.Label = *req.Label
	}
	if req.AssignedTo != nil {
		user.AssignedTo = req.AssignedTo
	}
	if req.Notes != nil {
		user.Notes = req.Notes
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, param *dto.DeleteUserParam) error {
	if err := s.validator.Validate(param); err != nil {
		return err
	}

	id, err := s.uuidPkg.Parse(param.ID)
	if err != nil {
		return errx.ErrUserNotFound.WithDetails(map[string]any{
			"id": param.ID,
		}).WithLocation("UserService.Delete").WithError(err)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
