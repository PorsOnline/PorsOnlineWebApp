package app

import (
	"github.com/porseOnline/config"
	notifPort "github.com/porseOnline/internal/notification/port"
	userPort "github.com/porseOnline/internal/user/port"
	surveyPort "github.com/porseOnline/internal/survey/port"
)

type App interface {
	UserService() userPort.Service
	NotifService() notifPort.Service
	SurveyService() surveyPort.Service
	Config() config.Config
}
