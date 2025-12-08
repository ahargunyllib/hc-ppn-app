package controller

import (
	"encoding/csv"
	"io"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
)

func (c *UserController) create(ctx *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	res, err := c.userSvc.Create(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

func (c *UserController) getByID(ctx *fiber.Ctx) error {
	var params dto.GetUserByIDParam
	if err := ctx.ParamsParser(&params); err != nil {
		return err
	}

	res, err := c.userSvc.GetByID(ctx.Context(), &params)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *UserController) list(ctx *fiber.Ctx) error {
	var query dto.GetUsersQuery
	if err := ctx.QueryParser(&query); err != nil {
		return err
	}

	res, err := c.userSvc.List(ctx.Context(), &query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *UserController) update(ctx *fiber.Ctx) error {
	var params dto.UpdateUserParam
	if err := ctx.ParamsParser(&params); err != nil {
		return err
	}

	var req dto.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := c.userSvc.Update(ctx.Context(), &params, &req); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusNoContent, nil)
}

func (c *UserController) delete(ctx *fiber.Ctx) error {
	var params dto.DeleteUserParam
	if err := ctx.ParamsParser(&params); err != nil {
		return err
	}

	if err := c.userSvc.Delete(ctx.Context(), &params); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusNoContent, nil)
}

func (c *UserController) getAllPhoneNumbers(ctx *fiber.Ctx) error {
	res, err := c.userSvc.GetAllPhoneNumbers(ctx.Context())
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *UserController) getMetrics(ctx *fiber.Ctx) error {
	res, err := c.userSvc.GetMetrics(ctx.Context())
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *UserController) importCSV(ctx *fiber.Ctx) error {
	// Get the uploaded file
	file, err := ctx.FormFile("file")
	if err != nil {
		return errx.ErrInvalidRequest.WithDetails(map[string]any{
			"error": "No file uploaded",
		}).WithLocation("UserController.importCSV").WithError(err)
	}

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return errx.ErrInternalServer.WithDetails(map[string]any{
			"error": "Failed to open file",
		}).WithLocation("UserController.importCSV").WithError(err)
	}
	defer fileContent.Close()

	// Parse CSV
	reader := csv.NewReader(fileContent)
	records, err := reader.ReadAll()
	if err != nil {
		if err == io.EOF {
			return errx.ErrInvalidRequest.WithDetails(map[string]any{
				"error": "Empty CSV file",
			}).WithLocation("UserController.importCSV")
		}
		return errx.ErrInvalidRequest.WithDetails(map[string]any{
			"error": "Failed to parse CSV file",
		}).WithLocation("UserController.importCSV").WithError(err)
	}

	// Validate CSV has at least header row
	if len(records) < 1 {
		return errx.ErrInvalidRequest.WithDetails(map[string]any{
			"error": "CSV file must have at least a header row",
		}).WithLocation("UserController.importCSV")
	}

	// Import users
	res, err := c.userSvc.ImportUsersFromCSV(ctx.Context(), records)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
