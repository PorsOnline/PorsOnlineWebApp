package app

import (
	"context"

	"github.com/porseOnline/config"
	codeVerificationPort "github.com/porseOnline/internal/codeVerification/port"
	notifPort "github.com/porseOnline/internal/notification/port"

	questionPort "github.com/porseOnline/internal/question/port"
	surveyPort "github.com/porseOnline/internal/survey/port"
	userPort "github.com/porseOnline/internal/user/port"
	votingPort "github.com/porseOnline/internal/voting/port"
)

type App interface {
	UserService() userPort.Service
	NotifService() notifPort.Service
	SurveyService() surveyPort.Service
	QuestionService() questionPort.Service
	PermissionService() userPort.PermissionService
	RoleService() userPort.RoleService
	VotingService() votingPort.Service

	surveyPort "github.com/porseOnline/internal/survey/port"
	userPort "github.com/porseOnline/internal/user/port"
	"gorm.io/gorm"
)

type App interface {
	UserService(ctx context.Context) userPort.Service
	NotifService(ctx context.Context) notifPort.Service
	SurveyService(ctx context.Context) surveyPort.Service
	CodeVerificationService(ctx context.Context) codeVerificationPort.Service
	DB() *gorm.DB

	Config() config.Config
}
