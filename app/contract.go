package app

import (
	"context"

	"github.com/porseOnline/config"
	codeVerificationPort "github.com/porseOnline/internal/codeVerification/port"
	notifPort "github.com/porseOnline/internal/notification/port"
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
