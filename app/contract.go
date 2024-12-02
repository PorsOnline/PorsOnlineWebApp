package app

import (
	"PorsOnlineWebApp/config"
	surveyPort "PorsOnlineWebApp/internal/survey/port"
)

type App interface {
	SurveyService() surveyPort.Service
	Config() config.Config
}