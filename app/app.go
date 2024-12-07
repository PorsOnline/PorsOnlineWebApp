package app

import (
	"context"

	"github.com/porseOnline/config"
	"github.com/porseOnline/internal/codeVerification"
	"github.com/porseOnline/internal/common"
	notifPort "github.com/porseOnline/internal/notification/port"
	"github.com/porseOnline/internal/survey"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"github.com/porseOnline/internal/user"
	userPort "github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/pkg/adapters/storage"
	"github.com/porseOnline/pkg/postgres"

	codeVerificationPort "github.com/porseOnline/internal/codeVerification/port"

	"github.com/go-co-op/gocron/v2"
	"github.com/porseOnline/internal/notification"
	appCtx "github.com/porseOnline/pkg/context"
	"gorm.io/gorm"
)

type app struct {
	db                *gorm.DB
	cfg               config.Config
	userService       userPort.Service
	notifService      notifPort.Service
	surveyService     surveyPort.Service
	codeVrfctnService codeVerificationPort.Service
}

func (a *app) DB() *gorm.DB {
	return a.db
}
func (a *app) UserService(ctx context.Context) userPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.userService == nil {
			a.userService = a.userServiceWithDB(a.db)
		}
		return a.userService
	}

	return a.userServiceWithDB(db)
}

func (a *app) userServiceWithDB(db *gorm.DB) userPort.Service {
	return user.NewService(storage.NewUserRepo(db))
}
func (a *app) NotifService(ctx context.Context) notifPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.notifService == nil {
			a.notifService = a.notifServiceWithDB(a.db)
		}
		return a.notifService
	}

	return a.notifServiceWithDB(db)
}
func (a *app) notifServiceWithDB(db *gorm.DB) notifPort.Service {
	return notification.NewService(storage.NewNotifRepo(db))
}

func (a *app) SurveyService(ctx context.Context) surveyPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.surveyService == nil {
			a.surveyService = a.surveyServiceWithDB(a.db)
		}
		return a.surveyService
	}

	return a.surveyServiceWithDB(db)
}

func (a *app) surveyServiceWithDB(db *gorm.DB) surveyPort.Service {
	return survey.NewService(storage.NewSurveyRepo(db))
}
func (a *app) codeVerificationServiceWithDB(db *gorm.DB) codeVerificationPort.Service {
	return codeVerification.NewService(
		a.userService, storage.NewOutboxRepo(db), storage.NewCodeVerificationRepo(db))
}

func (a *app) codeVerificationService(ctx context.Context) codeVerificationPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.codeVrfctnService == nil {
			a.codeVrfctnService = a.codeVerificationServiceWithDB(a.db)
		}
		return a.codeVrfctnService
	}

	return a.codeVerificationServiceWithDB(db)
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) setDB() error {
	db, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
		User:   a.cfg.DB.User,
		Pass:   a.cfg.DB.Password,
		Host:   a.cfg.DB.Host,
		Port:   a.cfg.DB.Port,
		DBName: a.cfg.DB.Database,
		Schema: a.cfg.DB.Schema,
	})
	postgres.GormMigrations(db)

	if err != nil {
		return err
	}

	a.db = db
	return nil
}

func NewApp(cfg config.Config) (App, error) {
	a := &app{
		cfg: cfg,
	}

	if err := a.setDB(); err != nil {
		return nil, err
	}

	a.userService = user.NewService(storage.NewUserRepo(a.db))

	a.notifService = notification.NewService(storage.NewNotifRepo(a.db))

	a.surveyService = survey.NewService(storage.NewSurveyRepo(a.db))
	a.codeVrfctnService = codeVerification.NewService(a.userService, storage.NewOutboxRepo(a.db), storage.NewCodeVerificationRepo(a.db))

	return a, a.registerOutboxHandlers()
}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}
func (a *app) registerOutboxHandlers() error {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	common.RegisterOutboxRunner(a.codeVerificationServiceWithDB(a.db), scheduler)

	scheduler.Start()

	return nil
}
