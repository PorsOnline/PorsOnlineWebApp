package app

import (
	"github.com/porseOnline/config"
	notifPort "github.com/porseOnline/internal/notification/port"
	userPort "github.com/porseOnline/internal/user/port"
)

type App interface {
	UserService() userPort.Service
	NotifService() notifPort.Service
	Config() config.Config
}
