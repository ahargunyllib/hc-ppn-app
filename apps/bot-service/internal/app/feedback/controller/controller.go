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

	// TODO: Add middleware for authentication and authorization
	feedbackRouter.Post("/", controller.create)
	feedbackRouter.Get("/", controller.list)
	feedbackRouter.Get("/metrics", controller.getMetrics)
	feedbackRouter.Get("/satisfaction-trend", controller.getSatisfactionTrend)
	feedbackRouter.Get("/:id", controller.getByID)
}
