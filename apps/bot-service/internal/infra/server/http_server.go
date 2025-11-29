package server

import (
	feedbackcontroller "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/controller"
	feedbackrepository "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/repository"
	feedbackservice "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/feedback/service"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/user/controller"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/user/repository"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/app/user/service"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/middlewares"
	errorhandler "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/helpers/http/error_handler"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/helpers/http/response"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/jwt"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/uuid"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HTTPServer interface {
	Start(port string)
	MountMiddlewares()
	MountRoutes(db *sqlx.DB)
	GetApp() *fiber.App
}

type httpServer struct {
	app *fiber.App
}

func NewHTTPServer() HTTPServer {
	config := fiber.Config{
		CaseSensitive: true,
		AppName:       "HC PPN Backend",
		ServerHeader:  "HC PPN Backend",
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		ErrorHandler:  errorhandler.ErrorHandler,
	}

	app := fiber.New(config)

	return &httpServer{
		app: app,
	}
}

func (s *httpServer) GetApp() *fiber.App {
	return s.app
}

func (s *httpServer) Start(port string) {
	if port[0] != ':' {
		port = ":" + port
	}

	err := s.app.Listen(port)

	if err != nil {
		log.Fatal(log.CustomLogInfo{
			"error": err.Error(),
		}, "[SERVER][Start] failed to start server")
	}
}

func (s *httpServer) MountMiddlewares() {
	s.app.Use(middlewares.LoggerConfig())
	s.app.Use(middlewares.Helmet())
	s.app.Use(middlewares.Compress())
	s.app.Use(middlewares.Cors())
	s.app.Use(middlewares.RecoverConfig())
}

func (s *httpServer) MountRoutes(db *sqlx.DB) {
	jwtService := jwt.Jwt
	validatorService := validator.Validator
	uuidService := uuid.UUID

	middleware := middlewares.NewMiddleware(jwtService)

	s.app.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "HC PPN Backend is running")
	})

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "HC PPN Backend is running")
	})

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, validatorService, uuidService)
	controller.InitUserController(v1, userService, middleware)

	feedbackRepo := feedbackrepository.NewFeedbackRepository(db)
	feedbackService := feedbackservice.NewFeedbackService(feedbackRepo, validatorService, uuidService)
	feedbackcontroller.InitFeedbackController(v1, feedbackService, middleware)

	s.app.Use(func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusNotFound, "Route not found")
	})
}
