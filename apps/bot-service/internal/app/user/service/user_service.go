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

func (s *UserService) ImportFromCSV(ctx context.Context, req *dto.ImportUsersFromCSVRequest) error {
	if err := s.validator.Validate(req); err != nil {
		return err
	}

	records, err := s.csvPkg.ParseFileHeader(req.File)
	if err != nil {
		return errx.ErrInternalServer.WithLocation("UserService.ImportFromCSV").WithError(err)
	}

	// Validate CSV is not empty
	if len(records) == 0 {
		return errx.ErrEmptyCSVFile.WithLocation("UserService.ImportFromCSV")
	}

	// Validate header row exists and has correct columns
	if len(records) < 2 {
		return errx.ErrCSVNoData.WithLocation("UserService.ImportFromCSV")
	}

	header := records[0]
	expectedColumns := 5
	if len(header) != expectedColumns {
		return errx.ErrInvalidCSVStructure.
			WithLocation("UserService.ImportFromCSV").
			WithDetails(map[string]any{
				"expected": expectedColumns,
				"got":      len(header),
			})
	}

	users := make([]entity.User, 0, len(records)-1)
	now := time.Now()

	for idx, record := range records {
		// skip header
		if idx == 0 {
			continue
		}

		// Validate record has correct number of columns
		if len(record) != expectedColumns {
			return errx.ErrInvalidCSVRow.
				WithLocation("UserService.ImportFromCSV").
				WithDetails(map[string]any{
					"row":      idx + 1,
					"expected": expectedColumns,
					"got":      len(record),
				})
		}

		phoneNumber := record[0]
		name := record[1]

		// Validate required fields
		if phoneNumber == "" {
			return errx.ErrMissingPhoneNumber.
				WithLocation("UserService.ImportFromCSV").
				WithDetails(map[string]any{"row": idx + 1})
		}
		if name == "" {
			return errx.ErrMissingName.
				WithLocation("UserService.ImportFromCSV").
				WithDetails(map[string]any{"row": idx + 1})
		}

		// Validate phone number format using the same validator
		phoneValidationErr := s.validator.Validate(&dto.CreateUserRequest{
			PhoneNumber: phoneNumber,
			Name:        name,
		})
		if phoneValidationErr != nil {
			return errx.ErrInvalidPhoneNumber.
				WithLocation("UserService.ImportFromCSV").
				WithDetails(map[string]any{
					"row":         idx + 1,
					"phoneNumber": phoneNumber,
				}).
				WithError(phoneValidationErr)
		}

		id, err := s.uuidPkg.NewV7()
		if err != nil {
			return errx.ErrInternalServer.WithLocation("UserService.ImportFromCSV").WithError(err)
		}

		var jobTitle *string
		if record[2] != "" {
			jobTitle = &record[2]
		}

		var gender *string
		if record[3] != "" {
			genderValue := record[3]
			// Validate gender value
			if genderValue != "male" && genderValue != "female" {
				return errx.ErrInvalidGender.
					WithLocation("UserService.ImportFromCSV").
					WithDetails(map[string]any{
						"row":    idx + 1,
						"gender": genderValue,
					})
			}
			gender = &genderValue
		}

		var dateOfBirth *time.Time
		if record[4] != "" {
			parsedDate, err := time.Parse(time.DateOnly, record[4])
			if err != nil {
				return errx.ErrInvalidDateFormat.WithDetails(map[string]any{
					"row":            idx + 1,
					"dateOfBirth":    record[4],
					"expectedFormat": "YYYY-MM-DD",
				}).WithLocation("UserService.ImportFromCSV").WithError(err)
			}
			dateOfBirth = &parsedDate
		}

		user := entity.User{
			ID:          id,
			PhoneNumber: phoneNumber,
			Name:        name,
			JobTitle:    jobTitle,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		users = append(users, user)
	}

	if err := s.userRepo.BulkCreate(ctx, users); err != nil {
		return err
	}

	return nil
}
