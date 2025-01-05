package app

import (
	"context"

	"github.com/porseOnline/config"
	codeVerificationPort "github.com/porseOnline/internal/codeVerification/port"
	notifPort "github.com/porseOnline/internal/notification/port"
	"gorm.io/gorm"

	questionPort "github.com/porseOnline/internal/question/port"
	surveyPort "github.com/porseOnline/internal/survey/port"
	userPort "github.com/porseOnline/internal/user/port"
	votingPort "github.com/porseOnline/internal/voting/port"
)

type App interface {
	DB() *gorm.DB

	Config(ctx context.Context) config.Config
	UserService(ctx context.Context) userPort.Service
	NotifService(ctx context.Context) notifPort.Service
	SurveyService(ctx context.Context) surveyPort.Service
	QuestionService(ctx context.Context) questionPort.Service
	PermissionService(ctx context.Context) userPort.PermissionService
	RoleService(ctx context.Context) userPort.RoleService
	VotingService(ctx context.Context) votingPort.Service
	CodeVerificationService(ctx context.Context) codeVerificationPort.Service
}
