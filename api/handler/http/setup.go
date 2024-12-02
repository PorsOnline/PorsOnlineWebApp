package http

import (
	"PorsOnlineWebApp/api/service"
	"PorsOnlineWebApp/app"
	"PorsOnlineWebApp/config"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run(appContainer app.App, config config.Config) error {
	app := fiber.New(fiber.Config{
		AppName:           "Survey v0.0.1",
		EnablePrintRoutes: true,
	})

	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        config.Server.RateLimitMaxAttempt,
		Expiration: time.Duration(config.Server.RatelimitTimePeriod) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("STOP` SENDING TOO MUCH REQUESTS")
		},
	}))
	surveyService := service.NewSurveyService(appContainer.SurveyService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyApi := app.Group("api/v1/survey")
	surveyApi.Post("", CreateSurvey(surveyService))
	surveyApi.Get(":uuid", GetSurvey(surveyService))
	return app.Listen(fmt.Sprintf("%v:%d",config.Server.IPAddress, config.Server.HttpPort))
}
