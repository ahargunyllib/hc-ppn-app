package controller

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/service"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type FeedbackController struct {
	feedbackSvc *service.FeedbackService
}

func InitFeedbackController(router fiber.Router, feedbackSvc *service.FeedbackService, middleware *middlewares.Middleware) {
	controller := &FeedbackController{
		feedbackSvc: feedbackSvc,
	}

	feedbackRouter := router.Group("/feedbacks")

	feedbackRouter.Post("/", middleware.RequireAuth(), controller.create)
	feedbackRouter.Get("/", middleware.RequireAuth(), controller.list)
	feedbackRouter.Get("/:id", middleware.RequireAuth(), controller.getByID)
}
