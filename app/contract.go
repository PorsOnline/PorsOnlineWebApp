package app

import (
	"github.com/porseOnline/config"
	notifPort "github.com/porseOnline/internal/notification/port"
	questionPort "github.com/porseOnline/internal/question/port"
	surveyPort "github.com/porseOnline/internal/survey/port"
	userPort "github.com/porseOnline/internal/user/port"
)

type App interface {
	UserService() userPort.Service
	NotifService() notifPort.Service
	SurveyService() surveyPort.Service
	QuestionService() questionPort.Service
	PermissionService() userPort.PermissionService
	RoleService() userPort.RoleService
	Config() config.Config
}
