package app

import (
	"github.com/porseOnline/config"
	"github.com/porseOnline/internal/question"
	questionPort "github.com/porseOnline/internal/question/port"
	"github.com/porseOnline/internal/survey"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"github.com/porseOnline/internal/user"
	userPort "github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/pkg/adapters/storage"
	"github.com/porseOnline/pkg/postgres"

	notifPort "github.com/porseOnline/internal/notification/port"

	"github.com/porseOnline/internal/notification"

	"gorm.io/gorm"
)

type app struct {
	db              *gorm.DB
	secretsDB       *gorm.DB
	cfg             config.Config
	userService     userPort.Service
	notifService    notifPort.Service
	surveyService   surveyPort.Service
	questionService questionPort.Service
}

func (a *app) UserService() userPort.Service {
	return a.userService
}

func (a *app) NotifService() notifPort.Service {
	return a.notifService
}

func (a *app) SurveyService() surveyPort.Service {
	return a.surveyService
}

func (a *app) QuestionService() questionPort.Service {
	return a.questionService
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

	secretDB, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
		User:   a.cfg.DB.User,
		Pass:   a.cfg.DB.Password,
		Host:   a.cfg.DB.Host,
		Port:   a.cfg.DB.Port,
		DBName: a.cfg.DB.SDatabase,
		Schema: a.cfg.DB.Schema,
	})
	postgres.GormSecretsMigration(secretDB)

	if err != nil {
		return err
	}

	a.secretsDB = secretDB

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

	a.questionService = question.NewService(storage.NewQuestionRepo(a.db), a.surveyService)

	return a, nil
}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}
