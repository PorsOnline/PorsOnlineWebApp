package app

import (
	"PorsOnlineWebApp/config"
	notifPort "PorsOnlineWebApp/internal/notification/port"
)

type App interface {
	NotifService() notifPort.Service
	Config() config.Config
}
