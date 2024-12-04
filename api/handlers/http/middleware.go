package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("X-Trace-ID", traceID)

		c.Locals("traceID", traceID)

		return c.Next()
	}
}
