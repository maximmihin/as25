package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"

	handlers "github.com/maximmihin/as25/internal/controllers/http"
	"github.com/maximmihin/as25/internal/logger"
)

func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {

		if err == nil {
			return nil
		}

		ctxAttr := logger.ExtractSlogAttrs(c.UserContext()) // TODO may be group them?
		for _, attr := range ctxAttr {
			slogfiber.AddCustomAttributes(c, attr)
		}

		fErr := new(fiber.Error) // TODO just ptr?
		if errors.As(err, &fErr) {
			return c.Status(fErr.Code).JSON(handlers.Error{
				Message: fErr.Message,
			})
		} else {
			return c.Status(500).JSON(handlers.Error{Message: "Internal Server Error"})
		}

	}
}
