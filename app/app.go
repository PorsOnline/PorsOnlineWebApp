package app

import (
	"PorsOnlineWebApp/config"
	surveyPort "PorsOnlineWebApp/internal/survey/port"
	"PorsOnlineWebApp/pkg/adapters/storage"
	"PorsOnlineWebApp/internal/survey"

	"gorm.io/gorm"
)

type app struct {
	db          *gorm.DB
	cfg         config.Config
	surveyServer surveyPort.Service
}

func (a *app) SurveyService() surveyPort.Service {
	return a.surveyServer
}

func (a *app) Config() config.Config {
	return a.cfg
}

// func (a *app) setDB() error {
// 	db, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
// 		User:   a.cfg.DB.User,
// 		Pass:   a.cfg.DB.Password,
// 		Host:   a.cfg.DB.Host,
// 		Port:   a.cfg.DB.Port,
// 		DBName: a.cfg.DB.Database,
// 		Schema: a.cfg.DB.Schema,
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	a.db = db
// 	return nil
// }

func NewApp(cfg config.Config) (App, error) {
	a := &app{
		cfg: cfg,
	}

	// if err := a.setDB(); err != nil {
	// 	return nil, err
	// }

	a.surveyServer = survey.NewSurveyService(storage.NewSurveyRepo(a.db))

	return a, nil
}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}