package middlewares

import (
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/env"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
)

func LoggerConfig() fiber.Handler {
	fields := []string{
		"referer",
		"ip",
		"url",
		"latency",
		"status",
		"method",
		"error",
	}

	if env.AppEnv.AppEnv != "production" {
		fields = append(fields, "body")
		fields = append(fields, "reqHeaders")
		fields = append(fields, "resHeaders")
	}

	logger := log.GetLogger()
	config := fiberzerolog.Config{
		Logger:          logger,
		FieldsSnakeCase: true,
		Fields:          fields,
		Messages: []string{
			"[LoggerMiddleware.LoggerConfig] Server error",
			"[LoggerMiddleware.LoggerConfig] Client error",
			"[LoggerMiddleware.LoggerConfig] Success",
		},
	}

	return fiberzerolog.New(config)
}
