package controller

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/topic/service"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type TopicController struct {
	topicSvc *service.TopicService
}

func InitTopicController(router fiber.Router, topicSvc *service.TopicService, middleware *middlewares.Middleware) {
	controller := &TopicController{
		topicSvc: topicSvc,
	}

	topicRouter := router.Group("/topics")

	// TODO: Add middleware for authentication and authorization
	topicRouter.Post("/bulk", controller.bulkCreate)
	topicRouter.Get("/hot", controller.getHotTopics)
}
