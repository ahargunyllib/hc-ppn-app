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

	var dateOfBirth *time.Time
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		parsedDate, err := time.Parse(time.DateOnly, *req.DateOfBirth)
		if err != nil {
			return nil, errx.ErrInvalidDateFormat.WithDetails(map[string]any{
				"req.DateOfBirth": *req.DateOfBirth,
			}).WithLocation("UserService.Create").WithError(err)
		}
		dateOfBirth = &parsedDate
	}

	user := &entity.User{
		ID:          id,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		JobTitle:    req.JobTitle,
		Gender:      req.Gender,
		DateOfBirth: dateOfBirth,
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

func (s *UserService) GetByPhoneNumber(ctx context.Context, param *dto.GetUserByPhoneNumberParam) (*dto.GetUserByPhoneNumberResponse, error) {
	// Validate param
	if err := s.validator.Validate(param); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByPhoneNumber(ctx, param.PhoneNumber)
	if err != nil {
		return nil, err
	}

	userRes := dto.ToUserResponse(user)

	res := &dto.GetUserByPhoneNumberResponse{
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
		Offset: (page - 1) * limit,
		Limit:  limit,
		Search: query.Search,
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
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.JobTitle != nil {
		user.JobTitle = req.JobTitle
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.DateOfBirth != nil {
		if *req.DateOfBirth == "" {
			user.DateOfBirth = nil
		} else {
			parsedDate, err := time.Parse(time.DateOnly, *req.DateOfBirth)
			if err != nil {
				return errx.ErrInvalidDateFormat.WithDetails(map[string]any{
					"req.DateOfBirth": *req.DateOfBirth,
				}).WithLocation("UserService.Update").WithError(err)
			}
			user.DateOfBirth = &parsedDate
		}
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

func (s *UserService) GetAllPhoneNumbers(ctx context.Context) (*dto.GetAllPhoneNumbersResponse, error) {
	phoneNumbers, err := s.userRepo.GetAllPhoneNumbers(ctx)
	if err != nil {
		return nil, err
	}

	res := &dto.GetAllPhoneNumbersResponse{
		PhoneNumbers: phoneNumbers,
	}

	return res, nil
}

func (s *UserService) GetMetrics(ctx context.Context) (*dto.GetUserMetricsResponse, error) {
	totalUsers, err := s.userRepo.GetTotalUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := &dto.GetUserMetricsResponse{
		TotalUsers: totalUsers,
	}

	return res, nil
}

func (s *UserService) ImportUsersFromCSV(ctx context.Context, records [][]string) (*dto.ImportUsersFromCSVResponse, error) {
	res := &dto.ImportUsersFromCSVResponse{
		Total:   len(records) - 1, // Exclude header row
		Success: 0,
		Failed:  0,
		Errors:  make([]dto.ImportUsersFromCSVError, 0),
	}

	// Collect valid users for bulk insert
	validUsers := make([]entity.User, 0, len(records)-1)

	// Skip header row and validate all records
	for i := 1; i < len(records); i++ {
		record := records[i]
		rowNumber := i + 1 // +1 because spreadsheets start at 1

		// Check if row has the correct number of columns (at least 2 for phoneNumber and name)
		if len(record) < 2 {
			res.Failed++
			res.Errors = append(res.Errors, dto.ImportUsersFromCSVError{
				Row:   rowNumber,
				Error: "Invalid row format: missing required columns",
			})
			continue
		}

		// Parse CSV columns (phoneNumber, name, jobTitle, gender, dateOfBirth)
		phoneNumber := record[0]
		name := record[1]

		var jobTitle *string
		if len(record) > 2 && record[2] != "" {
			jobTitle = &record[2]
		}

		var gender *string
		if len(record) > 3 && record[3] != "" {
			gender = &record[3]
		}

		var dateOfBirth *string
		if len(record) > 4 && record[4] != "" {
			dateOfBirth = &record[4]
		}

		// Create and validate user request
		userReq := &dto.CreateUserRequest{
			PhoneNumber: phoneNumber,
			Name:        name,
			JobTitle:    jobTitle,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
		}

		// Validate the data
		if err := s.validator.Validate(userReq); err != nil {
			res.Failed++
			res.Errors = append(res.Errors, dto.ImportUsersFromCSVError{
				Row:   rowNumber,
				Error: err.Error(),
			})
			continue
		}

		// Generate UUID
		id, err := s.uuidPkg.NewV7()
		if err != nil {
			res.Failed++
			res.Errors = append(res.Errors, dto.ImportUsersFromCSVError{
				Row:   rowNumber,
				Error: "Failed to generate UUID: " + err.Error(),
			})
			continue
		}

		// Parse date of birth if provided
		var parsedDateOfBirth *time.Time
		if dateOfBirth != nil && *dateOfBirth != "" {
			parsedDate, err := time.Parse(time.DateOnly, *dateOfBirth)
			if err != nil {
				res.Failed++
				res.Errors = append(res.Errors, dto.ImportUsersFromCSVError{
					Row:   rowNumber,
					Error: "Invalid date format: " + err.Error(),
				})
				continue
			}
			parsedDateOfBirth = &parsedDate
		}

		// Build user entity
		user := entity.User{
			ID:          id,
			PhoneNumber: userReq.PhoneNumber,
			Name:        userReq.Name,
			JobTitle:    userReq.JobTitle,
			Gender:      userReq.Gender,
			DateOfBirth: parsedDateOfBirth,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		validUsers = append(validUsers, user)
	}

	// Bulk insert all valid users
	if len(validUsers) > 0 {
		if err := s.userRepo.BulkCreate(ctx, validUsers); err != nil {
			return nil, errx.ErrInternalServer.WithDetails(map[string]any{
				"error": "Failed to bulk insert users",
			}).WithLocation("UserService.ImportUsersFromCSV").WithError(err)
		}
		res.Success = len(validUsers)
	}

	return res, nil
}
