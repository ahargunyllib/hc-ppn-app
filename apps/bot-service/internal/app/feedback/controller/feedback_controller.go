package controller

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
)

func (c *FeedbackController) create(ctx *fiber.Ctx) error {
	var req dto.CreateFeedbackRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	res, err := c.feedbackSvc.Create(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

func (c *FeedbackController) getByID(ctx *fiber.Ctx) error {
	var params dto.GetFeedbackByIDParam
	if err := ctx.ParamsParser(&params); err != nil {
		return err
	}

	res, err := c.feedbackSvc.GetByID(ctx.Context(), &params)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *FeedbackController) list(ctx *fiber.Ctx) error {
	var query dto.GetFeedbacksQuery
	if err := ctx.QueryParser(&query); err != nil {
		return err
	}

	res, err := c.feedbackSvc.List(ctx.Context(), &query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *FeedbackController) getMetrics(ctx *fiber.Ctx) error {
	res, err := c.feedbackSvc.GetMetrics(ctx.Context())
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
