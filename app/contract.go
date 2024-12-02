package app

import (

	"github.com/porseOnline/config"
	userPort "github.com/porseOnline/internal/user/port"
  notifPort "PorsOnlineWebApp/internal/notification/port"
)

type App interface {
	UserService() userPort.Service
  NotifService() notifPort.Service
  Config() config.Config
	
)

