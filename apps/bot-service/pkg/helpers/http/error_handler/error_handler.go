package errorhandler

import (
	"errors"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/helpers/http/response"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/log"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		return response.SendResponse(c, fiber.StatusUnprocessableEntity, map[string]any{
			"error":      fiber.Map{
				"message":    "Validation error",
				"error_code": "VALIDATION_ERROR",
				"details":    valErr,
			},
		})
	}

	var reqErr *errx.RequestError
	if errors.As(err, &reqErr) {
		log.Error(log.CustomLogInfo{
			"error_code": reqErr.ErrorCode,
			"location":   reqErr.Location,
			"details":    reqErr.Details,
			"error":      reqErr.Err,
		}, "[ErrorHandler] Request error")

		return response.SendResponse(c, reqErr.StatusCode, fiber.Map{
			"error": reqErr,
		})
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return response.SendResponse(c, fiberErr.Code, fiber.Map{
			"error": fiber.Map{
				"message": fiberErr.Message,
			},
		})
	}

	return response.SendResponse(c, fiber.StatusInternalServerError, errx.ErrInternalServer.WithError(err))
}
