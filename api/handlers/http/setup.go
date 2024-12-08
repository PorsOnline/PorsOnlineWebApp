package http

import (
	"fmt"
	"os"
	"time"

	"github.com/porseOnline/app"
	"github.com/porseOnline/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run(appContainer app.App, config config.ServerConfig) error {
	app := fiber.New(fiber.Config{
		AppName:           "Survey v0.0.1",
		EnablePrintRoutes: true,
	})

	app.Use(TraceMiddleware())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path} TraceID: ${locals:traceID}\n",
		Output: os.Stdout,
	}))
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        config.RateLimitMaxAttempt,
		Expiration: time.Duration(config.RatelimitTimePeriod) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			xForwardedFor := c.Get("x-forwarded-for")
			if xForwardedFor == "" {
				return c.IP()
			}
			return xForwardedFor
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("STOP SENDING TOO MUCH REQUESTS")
		},
	}))
	//surveyService := service.NewService(appContainer.SurveyService(), config.Secret, config.AuthExpMinute, config.AuthRefreshMinute)

	// userService := service.NewUserService(appContainer.UserService(),
	// 	config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	// notifService := service.NewNotificationSerivce(appContainer.NotifService(), config.Secret, config.AuthExpMinute, config.AuthRefreshMinute)
	// api := app.Group("/api/v1/")

	api := app.Group("/api/v1")
	surveyApi := api.Group("/survey")
	userApi := api.Group("/user")
	notifApi := api.Group("/notif")

	registerAuthAPI(appContainer, config, userApi, surveyApi, notifApi)
	return app.Listen(fmt.Sprintf(":%d", config.HttpPort))
}
func registerAuthAPI(appContainer app.App, cfg config.ServerConfig, userRouter fiber.Router, surveyRouter fiber.Router, notifRouter fiber.Router) {
	userSvcGetter := userServiceGetter(appContainer, cfg)
	surveySvcGetter := SurveyServiceGetter(appContainer, cfg)
	notifSvcGetter := NotificationServiceGetter(appContainer, cfg)
	//user
	userRouter.Post("/sign-up", SignUp(userSvcGetter))
	userRouter.Post("/sign-in", SignIn(userSvcGetter))
	userRouter.Post("/sign-up-code-verification", SignUpCodeVerification(userSvcGetter))
	userRouter.Get("/users/:id", GetUserByID(userSvcGetter))
	//notif
	notifRouter.Post("/send_message", SendMessage(notifSvcGetter))
	notifRouter.Get("/unread-messages/:user_id", GetUnreadMessages(notifSvcGetter))
	//survey

	surveyRouter.Post("", CreateSurvey(surveySvcGetter))
	surveyRouter.Get(":uuid", GetSurvey(surveySvcGetter))
	surveyRouter.Put(":uuid", UpdateSurvey(surveySvcGetter))
	surveyRouter.Post("cancel/:uuid", CancelSurvey(surveySvcGetter))
	surveyRouter.Delete(":uuid", DeleteSurvey(surveySvcGetter))
	surveyRouter.Get("", GetAllSurveys(surveySvcGetter))
}
