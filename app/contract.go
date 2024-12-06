package app

import (
	"github.com/porseOnline/config"
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
	VotingService() votingPort.Service
	Config() config.Config
}
