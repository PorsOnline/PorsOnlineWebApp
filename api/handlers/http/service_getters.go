package http

import (
	"context"

	"github.com/porseOnline/api/service"
	"github.com/porseOnline/app"
	"github.com/porseOnline/config"
)

// user service transient instance handler
func userServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.UserService] {
	return func(ctx context.Context) *service.UserService {
		return service.NewUserService(appContainer.UserService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute, appContainer.CodeVerificationService(ctx))
	}
}
func SurveyServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.SurveyService] {
	return func(ctx context.Context) *service.SurveyService {
		return service.NewService(appContainer.SurveyService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}
func NotificationServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.NotificationService] {
	return func(ctx context.Context) *service.NotificationService {
		return service.NewNotificationSerivce(appContainer.NotifService(ctx), cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}
