package controller

import (
	"net/http"
	"slices"

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
	var req dto.ImportUsersFromCSVRequest
	var err error

	req.File, err = ctx.FormFile("file")
	if err != nil {
		return err
	}

	// check mime type
	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	const mimeDetectionBufferSize = 512
	buffer := make([]byte, mimeDetectionBufferSize)
	n, err := file.Read(buffer)
	if err != nil {
		return err
	}
	// reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	mimeType := http.DetectContentType(buffer[:n])
	allowedMimeTypes := []string{"text/csv", "application/vnd.ms-excel", "text/plain"}
	if !slices.Contains(allowedMimeTypes, mimeType) {
		return errx.ErrInvalidFileExtension.WithLocation("userController.ImportCSV").WithDetails(map[string]any{
			"expected": "text/csv or application/vnd.ms-excel",
			"got":      mimeType,
		})
	}

	// check file size (max 5MB)
	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if req.File.Size > maxFileSize {
		return errx.ErrFileSizeLimitExceeded.WithLocation("userController.ImportCSV").WithDetails(map[string]any{
			"maxSize": maxFileSize,
			"got":     req.File.Size,
		})
	}

	err = c.userSvc.ImportFromCSV(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}
