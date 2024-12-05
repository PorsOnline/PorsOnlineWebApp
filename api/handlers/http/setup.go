package http

import (
	"fmt"
	"os"
	"time"

	"github.com/porseOnline/api/service"
	"github.com/porseOnline/app"
	"github.com/porseOnline/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run(appContainer app.App, config config.Config) error {
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
		Max:        config.Server.RateLimitMaxAttempt,
		Expiration: time.Duration(config.Server.RatelimitTimePeriod) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("STOP` SENDING TOO MUCH REQUESTS")
		},
	}))
	surveyService := service.NewService(appContainer.SurveyService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyApi := app.Group("api/v1/survey")
	surveyApi.Post("", CreateSurvey(surveyService))
	surveyApi.Get(":uuid", GetSurvey(surveyService))
	surveyApi.Put("", UpdateSurvey(surveyService))
	surveyApi.Post("cancel/:uuid", CancelSurvey(surveyService))
	surveyApi.Delete(":uuid", DeleteSurvey(surveyService))
	surveyApi.Get("", GetAllSurveys(surveyService))
	userService := service.NewUserService(appContainer.UserService(),
		config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)

	api := app.Group("/api/v1")
	api.Post("/sign-up", SignUp(userService))
	api.Post("/sign-in", SignIn(userService))
	api.Post("/sign-up-code-verification", SignUpCodeVerification(userService))
	api.Post("/user/update", Update(userService))

	api.Get("/users/:id", GetUserByID(userService))
	notifService := service.NewNotificationSerivce(appContainer.NotifService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	api.Post("/send_message", SendMessage(notifService))
	api.Get("/unread-messages/:user_id", GetUnreadMessages(notifService))

	questionService := service.NewQuestionService(appContainer.QuestionService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyApi.Post("/question", CreateQuestion(questionService))
	surveyApi.Delete("/question/:id", DeleteQuestion(questionService))
	surveyApi.Put("/question", UpdateQuestion(questionService))

	return app.Listen(fmt.Sprintf(":%d", config.Server.HttpPort))
}
