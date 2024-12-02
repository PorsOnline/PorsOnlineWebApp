package app

import (
	"github.com/porseOnline/config"
	userPort "github.com/porseOnline/internal/user/port"
)

type App interface {
	UserService() userPort.Service
	Config() config.Config
}
