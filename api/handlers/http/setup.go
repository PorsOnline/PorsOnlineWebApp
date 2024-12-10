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
//<<<<<<< dev
	permissionService := service.NewPermissionService(appContainer.PermissionService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyService := service.NewService(appContainer.SurveyService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	questionService := service.NewQuestionService(appContainer.QuestionService(), config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	surveyApi := app.Group("api/v1/survey")




	

	//surveyService := service.NewService(appContainer.SurveyService(), config.Secret, config.AuthExpMinute, config.AuthRefreshMinute)

	// userService := service.NewUserService(appContainer.UserService(),
	// 	config.Server.Secret, config.Server.AuthExpMinute, config.Server.AuthRefreshMinute)
	// notifService := service.NewNotificationSerivce(appContainer.NotifService(), config.Secret, config.AuthExpMinute, config.AuthRefreshMinute)
	// api := app.Group("/api/v1/")

	api := app.Group("/api/v1")
	surveyApi := api.Group("/survey")
	userApi := api.Group("/user")
	notifApi := api.Group("/notif")
//>>>>>>> feat/validation-code-and-its-time

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
	userRouter.Put("/user/update", Update(userService))
	userRouter.Delete("/user/:id", PermissionMiddleware(permissionService), DeleteByID(userService))
	//notif
	notifRouter.Post("/send_message", SendMessage(notifSvcGetter))
	notifRouter.Get("/unread-messages/:user_id", GetUnreadMessages(notifSvcGetter))
	//survey
//<<<<<<< dev


	surveyApi.Use(newAuthMiddleware([]byte(config.Server.Secret)))
	surveyApi.Post(":surveyID/question", PermissionMiddleware(permissionService), CreateQuestion(questionService))
	surveyApi.Delete(":surveyID/question/:id", PermissionMiddleware(permissionService), DeleteQuestion(questionService))
	surveyApi.Put(":surveyID/question", PermissionMiddleware(permissionService), UpdateQuestion(questionService))
	surveyApi.Get(":surveyID/question/get-next", PermissionMiddleware(permissionService), UpdateQuestion(questionService))
	surveyApi.Post("", CreateSurvey(surveyService))
	surveyApi.Get(":surveyID", PermissionMiddleware(permissionService),GetSurvey(surveyService))
	surveyApi.Put(":surveyID", PermissionMiddleware(permissionService),UpdateSurvey(surveyService))
	surveyApi.Post("cancel/:surveyID", PermissionMiddleware(permissionService),CancelSurvey(surveyService))
	surveyApi.Delete(":surveyID", PermissionMiddleware(permissionService),DeleteSurvey(surveyService))
	surveyApi.Get("", PermissionMiddleware(permissionService),GetAllSurveys(surveyService))

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
//=======
	surveyRouter.Post("", CreateSurvey(surveySvcGetter))
	surveyRouter.Get(":uuid", GetSurvey(surveySvcGetter))
	surveyRouter.Put(":uuid", UpdateSurvey(surveySvcGetter))
	surveyRouter.Post("cancel/:uuid", CancelSurvey(surveySvcGetter))
	surveyRouter.Delete(":uuid", DeleteSurvey(surveySvcGetter))
	surveyRouter.Get("", GetAllSurveys(surveySvcGetter))
//>>>>>>> feat/validation-code-and-its-time
}
