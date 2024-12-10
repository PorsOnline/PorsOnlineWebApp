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

	permissionService := service.NewPermissionService(appContainer.PermissionService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyService := service.NewService(appContainer.SurveyService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	questionService := service.NewQuestionService(appContainer.QuestionService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)

	api := app.Group("/api/v1")
	surveyApi := api.Group("/survey")
	userApi := api.Group("/user")
	notifApi := api.Group("/notif")
	votingApi := app.Group("/vote")
	roleApi := app.Group("/role")
	permissionApi := app.Group("/permission")
	roleService := service.NewRoleService(appContainer.RoleService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	votingService := service.NewVotingService(appContainer.VotingService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)

	registerAuthAPI(appContainer, config, userApi, surveyApi, notifApi, votingApi, permissionApi, roleApi)
	return app.Listen(fmt.Sprintf(":%d", config.HttpPort))
}
func registerAuthAPI(appContainer app.App, cfg config.ServerConfig, userRouter fiber.Router, surveyRouter fiber.Router, notifRouter fiber.Router, votingRouter fiber.Router, permissionRouter fiber.Router, roleRouter fiber.Router) {
	userSvcGetter := userServiceGetter(appContainer, cfg)
	surveySvcGetter := SurveyServiceGetter(appContainer, cfg)
	notifSvcGetter := NotificationServiceGetter(appContainer, cfg)
	voteSvcGetter := voteSvcGetter(appContainer, cfg)
	roleSvcGetter := roleSvcGetter(appContainer, cfg)
	permissionSvcGetter := permissionSvcGetter(appContainer, cfg)
	questionSvcGetter := questionSvcGetter(appContainer, cfg)
	//user
	userRouter.Post("/sign-up", SignUp(userSvcGetter))
	userRouter.Post("/sign-in", SignIn(userSvcGetter))
	userRouter.Post("/sign-up-code-verification", SignUpCodeVerification(userSvcGetter))
	userRouter.Get("/users/:id", GetUserByID(userSvcGetter))
	userRouter.Put("/user/update", Update(userService))
	userRouter.Delete("/user/:id", PermissionMiddleware(permissionService), DeleteByID(userService))
	//notif
	notifRouter.Post("/send_message", SendMessage(notifSvcGetter))
	notifRouter.Get("/unread-messages/:user_id", GetUnreadMessages(notifSvcGetter))
	//survey
	surveyRouter.Use(newAuthMiddleware([]byte(config.Server.Secret)))
	surveyRouter.Post("", CreateSurvey(surveySvcGetter))
	surveyRouter.Post(":surveyID/question", PermissionMiddleware(permissionService), CreateQuestion(questionService))
	surveyRouter.Delete(":surveyID/question/:id", PermissionMiddleware(permissionService), DeleteQuestion(questionService))
	surveyRouter.Put(":surveyID/question", PermissionMiddleware(permissionService), UpdateQuestion(questionService))
	surveyRouter.Get(":surveyID/question/get-next", PermissionMiddleware(permissionService), UpdateQuestion(questionService))
	surveyRouter.Post("", CreateSurvey(surveySvcGetter))
	surveyRouter.Get(":surveyID", PermissionMiddleware(permissionService), GetSurvey(surveySvcGetter))
	surveyRouter.Put(":surveyID", PermissionMiddleware(permissionService), UpdateSurvey(surveySvcGetter))
	surveyRouter.Post("cancel/:surveyID", PermissionMiddleware(permissionService), CancelSurvey(surveySvcGetter))
	surveyRouter.Delete(":surveyID", PermissionMiddleware(permissionService), DeleteSurvey(surveySvcGetter))
	surveyRouter.Get("", PermissionMiddleware(permissionService), GetAllSurveys(surveySvcGetter))

	roleRouter.Post("", CreateRole(roleSvcGetter))
	roleRouter.Get(":id", GetRole(roleSvcGetter))
	roleRouter.Put("", UpdateRole(roleSvcGetter))
	roleRouter.Delete(":id", DeleteRole(roleSvcGetter))
	roleRouter.Patch(":roleId/assign/:userId", AssignRoleToUser(roleSvcGetter))

	permissionRouter.Post("", CreatePermission(permissionSvcGetter))
	permissionRouter.Get(":id", GetUserPermissions(permissionSvcGetter))
	permissionRouter.Get(":id", GetPermissionByID(permissionSvcGetter))
	permissionRouter.Put("", UpdatePermission(permissionSvcGetter))
	permissionRouter.Delete(":id", DeletePermission(permissionSvcGetter))
	permissionRouter.Patch(":permissionId/assign/:userId", AssignPermissionToUser(permissionSvcGetter))

	votingRouter.Post("", Vote(votingService))

	return app.Listen(fmt.Sprintf(":%d", config.Server.HttpPort))

}
