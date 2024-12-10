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

func surveyServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.SurveyService] {
	return func(ctx context.Context) *service.SurveyService {
		return service.NewService(appContainer.SurveyService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}
func notificationServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.NotificationService] {
	return func(ctx context.Context) *service.NotificationService {
		return service.NewNotificationSerivce(appContainer.NotifService(ctx), cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}
func roleServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.RoleService] {
	return func(ctx context.Context) *service.RoleService {
		return service.NewRoleService(appContainer.RoleService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}
func permissionServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.PermissionService] {
	return func(ctx context.Context) *service.PermissionService {
		return service.NewPermissionService(appContainer.PermissionService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}
func votingServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.VoteService] {
	return func(ctx context.Context) *service.VoteService {
		return service.NewVotingService(appContainer.VotingService(ctx), cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}
func questionSvcGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.QuestionService] {
	return func(ctx context.Context) *service.QuestionService {
		return service.NewQuestionService(appContainer.QuestionService(ctx), cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}
