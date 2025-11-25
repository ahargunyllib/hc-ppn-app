package controller

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/user/service"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userSvc *service.UserService
}

func InitUserController(router fiber.Router, userSvc *service.UserService, middleware *middlewares.Middleware) {
	controller := &UserController{
		userSvc: userSvc,
	}

	userRouter := router.Group("/users")

	userRouter.Post("/", middleware.RequireAuth(), controller.create)
	userRouter.Get("/", middleware.RequireAuth(), controller.list)
	userRouter.Patch("/:id", middleware.RequireAuth(), controller.update)
	userRouter.Delete("/:id", middleware.RequireAuth(), controller.delete)
	userRouter.Get("/:id", middleware.RequireAuth(), controller.getByID)
}
