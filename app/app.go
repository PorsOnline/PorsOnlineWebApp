package app

import (
	"PorsOnlineWebApp/config"
	notifPort "PorsOnlineWebApp/internal/notification/port"
	"PorsOnlineWebApp/pkg/adapters/storage"
	"PorsOnlineWebApp/internal/notification"

	"gorm.io/gorm"
)

type app struct {
	db          *gorm.DB
	cfg         config.Config
	notifServer notifPort.Service
}

func (a *app) NotifService() notifPort.Service {
	return a.notifServer
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

	a.notifServer = notification.NewService(storage.NewNotifRepo(a.db))

	return a, nil
}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}
