package http

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/porseOnline/app"
	"github.com/porseOnline/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/porseOnline/api/service"
)

func Run(appContainer app.App, config config.ServerConfig) error {
	app := fiber.New(fiber.Config{
		AppName:           "Survey v0.0.1",
	})
	app.Use(func(c *fiber.Ctx) error {
		permissionService := appContainer.PermissionService
		c.Locals("permissionService", permissionService)
		return c.Next()
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

	api := app.Group("/api/v1")

	permissionService := service.NewPermissionService(appContainer.PermissionService(context.Background()), config.Secret, config.AuthExpMinute, config.AuthRefreshMinute)
	registerAPI(appContainer, config, permissionService, api)
  
  	certFile := "/app/server.crt"
	keyFile := "/app/server.key"

	return app.ListenTLS(fmt.Sprintf(":%d", config.HttpPort), certFile, keyFile)
}
func registerAPI(appContainer app.App, cfg config.ServerConfig, permissionService *service.PermissionService, api fiber.Router) {
	surveyRouter := api.Group("/survey")
	userRouter := api.Group("/user")
	notifRouter := api.Group("/notif")
	votingRouter := api.Group("/vote")
	roleRouter := api.Group("/role")
	permissionRouter := api.Group("/permission")
	userSvcGetter := userServiceGetter(appContainer, cfg)
	surveySvcGetter := surveyServiceGetter(appContainer, cfg)
	notifSvcGetter := notificationServiceGetter(appContainer, cfg)
	voteSvcGetter := votingServiceGetter(appContainer, cfg)
	roleSvcGetter := roleServiceGetter(appContainer, cfg)
	permissionSvcGetter := permissionServiceGetter(appContainer, cfg)
	questionSvcGetter := questionSvcGetter(appContainer, cfg)
	//user
	userRouter.Post("/sign-up", SignUp(userSvcGetter))
	userRouter.Post("/sign-in", SignIn(userSvcGetter))
	userRouter.Post("/sign-up-code-verification", SignUpCodeVerification(userSvcGetter))
	userRouter.Get("/users/:id", GetUserByID(userSvcGetter))
	userRouter.Put("/user/update", Update(userSvcGetter))
	userRouter.Delete("/user/:id", PermissionMiddleware(permissionService), DeleteByID(userSvcGetter))
	//notif
	notifRouter.Post("/send_message", SendMessage(notifSvcGetter))
	notifRouter.Get("/unread-messages/:user_id", GetUnreadMessages(notifSvcGetter))
	//survey
	surveyRouter.Use(newAuthMiddleware([]byte(cfg.Secret)))
	surveyRouter.Post("", CreateSurvey(surveySvcGetter))
	surveyRouter.Post(":surveyID/question", PermissionMiddleware(permissionService), CreateQuestion(questionSvcGetter))
	surveyRouter.Delete(":surveyID/question/:id", PermissionMiddleware(permissionService), DeleteQuestion(questionSvcGetter))
	surveyRouter.Put(":surveyID/question", PermissionMiddleware(permissionService), UpdateQuestion(questionSvcGetter))
	surveyRouter.Get(":surveyID/question/get-next", PermissionMiddleware(permissionService), UpdateQuestion(questionSvcGetter))
	surveyRouter.Post("", CreateSurvey(surveySvcGetter))
	surveyRouter.Get(":surveyID", PermissionMiddleware(permissionService), GetSurvey(surveySvcGetter))
	surveyRouter.Put(":surveyID", PermissionMiddleware(permissionService), UpdateSurvey(surveySvcGetter))
	surveyRouter.Post("cancel/:surveyID", PermissionMiddleware(permissionService), CancelSurvey(surveySvcGetter))
	surveyRouter.Delete(":surveyID", PermissionMiddleware(permissionService), DeleteSurvey(surveySvcGetter))
	surveyRouter.Get("", PermissionMiddleware(permissionService), GetAllSurveys(surveySvcGetter))
	//role
	roleRouter.Post("", CreateRole(roleSvcGetter))
	roleRouter.Get(":id", GetRole(roleSvcGetter))
	roleRouter.Put("", UpdateRole(roleSvcGetter))
	roleRouter.Delete(":id", DeleteRole(roleSvcGetter))
	roleRouter.Patch(":roleId/assign/:userId", AssignRoleToUser(roleSvcGetter))
	//permission
	permissionRouter.Post("", CreatePermission(permissionSvcGetter))
	permissionRouter.Get(":id", GetUserPermissions(permissionSvcGetter))
	permissionRouter.Get(":id", GetPermissionByID(permissionSvcGetter))
	permissionRouter.Put("", UpdatePermission(permissionSvcGetter))
	permissionRouter.Delete(":id", DeletePermission(permissionSvcGetter))
	permissionRouter.Patch(":permissionId/assign/:userId", AssignPermissionToUser(permissionSvcGetter))
	//vote
	votingRouter.Post("", Vote(voteSvcGetter))

}