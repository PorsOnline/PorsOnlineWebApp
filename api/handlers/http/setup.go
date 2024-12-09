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
	permissionService := service.NewPermissionService(appContainer.PermissionService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyService := service.NewService(appContainer.SurveyService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyApi := app.Group("api/v1/survey")
	surveyApi.Use(newAuthMiddleware([]byte(config.Server.Secret)))
	surveyApi.Use(PermissionMiddleware(permissionService))
	surveyApi.Post("", CreateSurvey(surveyService))
	surveyApi.Get(":surveyID", GetSurvey(surveyService))
	surveyApi.Put("/:surveyID", UpdateSurvey(surveyService))
	surveyApi.Post("cancel/:surveyID", CancelSurvey(surveyService))
	surveyApi.Delete(":surveyID", DeleteSurvey(surveyService))
	surveyApi.Get("", GetAllSurveys(surveyService))
	userService := service.NewUserService(appContainer.UserService(),
		config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)

	api := app.Group("/api/v1")
	api.Post("/sign-up", SignUp(userService))
	api.Post("/sign-in", SignIn(userService))
	api.Post("/sign-up-code-verification", SignUpCodeVerification(userService))
	api.Put("/user/update", Update(userService))
	api.Delete("/user/:id", PermissionMiddleware(permissionService), DeleteByID(userService))

	api.Get("/users/:id", GetUserByID(userService))
	notifService := service.NewNotificationSerivce(appContainer.NotifService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	api.Post("/send_message", SendMessage(notifService))
	api.Get("/unread-messages/:user_id", GetUnreadMessages(notifService))

	questionService := service.NewQuestionService(appContainer.QuestionService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyApi.Post(":surveyID/question", PermissionMiddleware(permissionService), CreateQuestion(questionService))
	surveyApi.Delete(":surveyID/question/:id", PermissionMiddleware(permissionService), DeleteQuestion(questionService))
	surveyApi.Put(":surveyID/question", PermissionMiddleware(permissionService), UpdateQuestion(questionService))

	roleService := service.NewRoleService(appContainer.RoleService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	roleApi := app.Group("api/v1")
	roleApi.Post("/role", CreateRole(roleService))
	roleApi.Get("/role/:id", GetRole(roleService))
	roleApi.Put("/role", UpdateRole(roleService))
	roleApi.Delete("/role/:id", DeleteRole(roleService))
	roleApi.Patch("/role/:roleId/assign/:userId", AssignRoleToUser(roleService))

	permissionApi := app.Group("api/v1")
	permissionApi.Post("/permission", CreatePermission(permissionService))
	permissionApi.Get("/permissions/:id", GetUserPermissions(permissionService))
	permissionApi.Get("/permission/:id", GetPermissionByID(permissionService))
	permissionApi.Put("/permission", UpdatePermission(permissionService))
	permissionApi.Delete("/permission/:id", DeletePermission(permissionService))
	permissionApi.Patch("/permission/:permissionId/assign/:userId", AssignPermissionToUser(permissionService))

	votingApi := app.Group("api/v1/vote")
	votingService := service.NewVotingService(appContainer.VotingService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	votingApi.Post("", Vote(votingService))

	return app.Listen(fmt.Sprintf(":%d", config.Server.HttpPort))
}
