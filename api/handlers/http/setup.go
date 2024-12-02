package http

import (
	"fmt"

	"github.com/porseOnline/api/service"
	"github.com/porseOnline/app"
	"github.com/porseOnline/config"

	"github.com/gofiber/fiber/v2"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {
	router := fiber.New()

	api := router.Group("/api/v1")

	userService := service.NewUserService(appContainer.UserService(),
		cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)

	api.Post("/sign-up", SignUp(userService))
	api.Post("/sign-up-code-verification", SignUpCodeVerification(userService))

	api.Get("/users/:id", GetUserByID(userService))

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}
