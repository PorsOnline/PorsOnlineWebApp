package app

import (
	"context"

	"github.com/porseOnline/config"
	"github.com/porseOnline/internal/question"
	questionPort "github.com/porseOnline/internal/question/port"
	"github.com/porseOnline/internal/survey"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"github.com/porseOnline/internal/user"
	"github.com/porseOnline/internal/user/domain"
	userPort "github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/internal/voting"
	votingPort "github.com/porseOnline/internal/voting/port"
	"github.com/porseOnline/pkg/adapters/storage"
	"github.com/porseOnline/pkg/postgres"

	notifPort "github.com/porseOnline/internal/notification/port"

	"github.com/porseOnline/internal/notification"

	"gorm.io/gorm"
)

type app struct {
	db                *gorm.DB
	secretsDB         *gorm.DB
	cfg               config.Config
	userService       userPort.Service
	notifService      notifPort.Service
	surveyService     surveyPort.Service
	questionService   questionPort.Service
	votingService     votingPort.Service
	roleService       userPort.RoleService
	permissionService userPort.PermissionService
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

func (a *app) RoleService() userPort.RoleService {
	return a.roleService
}

func (a *app) PermissionService() userPort.PermissionService {
	return a.permissionService
}

func (a *app) VotingService() votingPort.Service {
	return a.votingService
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

	a.secretsDB = nil //secretDB

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

	a.roleService = user.NewRoleService(storage.NewRoleRepo(a.db))

	a.permissionService = user.NewPermissionService(storage.NewPermissionRepo(a.db), a.surveyService)

	a.notifService = notification.NewService(storage.NewNotifRepo(a.db))

	a.surveyService = survey.NewService(storage.NewSurveyRepo(a.db), a.permissionService)

	a.questionService = question.NewService(storage.NewQuestionRepo(a.db), a.surveyService)

	a.votingService = voting.NewVotingService(storage.NewVotingRepo(a.db, a.secretsDB))

	a.permissionService.SeedPermissions(context.Background(), generatePermissions())

	return a, nil
}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}

func generatePermissions() []domain.Permission {
	permissions := []domain.Permission{
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey/:uuid", Scope: "read"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey", Scope: "update"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey/cancel/:uuid", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey/:uuid", Scope: "delete"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey", Scope: "read"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/user", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/user/update", Scope: "update"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/user/:id", Scope: "delete"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/user/:id", Scope: "read"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/send_message", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/unread-messages/:user_id", Scope: "read"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey/:id/question", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey/:id/question/:id", Scope: "delete"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey/:id/question", Scope: "update"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/survey/:id/question/get-next", Scope: "read"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/role", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/role/:id", Scope: "read"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/role", Scope: "update"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/role/:id", Scope: "delete"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/role/:roleId/assign/:userId", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/permission", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/permissions/:id", Scope: "read"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/permission/:id", Scope: "read"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/permission", Scope: "update"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/permission/:id", Scope: "delete"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/permission/:userId/validate", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/permission/:permissionId/assign/:userId", Scope: "create"},
		{Policy: domain.PolicyUnknown, Resource: "/api/v1/vote", Scope: "create"},
	}
	return permissions
}