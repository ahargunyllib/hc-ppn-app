package controller

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/dto"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
)

func (c *TopicController) bulkCreate(ctx *fiber.Ctx) error {
	var req dto.BulkCreateTopicsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	err := c.topicSvc.BulkCreate(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, nil)
}

func (c *TopicController) getHotTopics(ctx *fiber.Ctx) error {
	res, err := c.topicSvc.GetHotTopics(ctx.Context())
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
