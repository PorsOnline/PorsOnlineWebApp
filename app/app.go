package app

import (
	"context"

	"github.com/porseOnline/config"

	"github.com/porseOnline/internal/question"
	questionPort "github.com/porseOnline/internal/question/port"

	"github.com/porseOnline/internal/codeVerification"
	"github.com/porseOnline/internal/common"
	notifPort "github.com/porseOnline/internal/notification/port"

	"github.com/porseOnline/internal/survey"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"github.com/porseOnline/internal/user"
	"github.com/porseOnline/internal/user/domain"
	userPort "github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/internal/voting"
	votingPort "github.com/porseOnline/internal/voting/port"
	"github.com/porseOnline/pkg/adapters/storage"
	"github.com/porseOnline/pkg/postgres"

	codeVerificationPort "github.com/porseOnline/internal/codeVerification/port"

	"github.com/go-co-op/gocron/v2"
	"github.com/porseOnline/internal/notification"
	appCtx "github.com/porseOnline/pkg/context"
	"gorm.io/gorm"
)

type app struct {
	db *gorm.DB

	secretsDB *gorm.DB

	cfg               config.Config
	userService       userPort.Service
	notifService      notifPort.Service
	surveyService     surveyPort.Service
	questionService   questionPort.Service
	votingService     votingPort.Service
	roleService       userPort.RoleService
	permissionService userPort.PermissionService
	codeVrfctnService codeVerificationPort.Service
}

// CodeVerificationService implements App.

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
	return survey.NewService(storage.NewSurveyRepo(db), a.permissionService)
}

func (a *app) codeVerificationServiceWithDB(db *gorm.DB) codeVerificationPort.Service {
	return codeVerification.NewService(
		a.userService, storage.NewOutboxRepo(db), storage.NewCodeVerificationRepo(db))
}

func (a *app) CodeVerificationService(ctx context.Context) codeVerificationPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.codeVrfctnService == nil {
			a.codeVrfctnService = a.codeVerificationServiceWithDB(a.db)
		}
		return a.codeVrfctnService
	}

	return a.codeVerificationServiceWithDB(db)
}

func (a *app) QuestionService(ctx context.Context) questionPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.questionService == nil {
			a.questionService = a.questionServiceWithDB(a.db)
		}
		return a.questionService
	}

	return a.questionServiceWithDB(db)
}

func (a *app) questionServiceWithDB(db *gorm.DB) questionPort.Service {
	return question.NewService(storage.NewQuestionRepo(db), a.surveyService)
}
func (a *app) RoleService(ctx context.Context) userPort.RoleService {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.roleService == nil {
			a.roleService = a.roleServiceWithDB(a.db)
		}
		return a.roleService
	}

	return a.roleServiceWithDB(db)
}

func (a *app) roleServiceWithDB(db *gorm.DB) userPort.RoleService {
	return user.NewRoleService(storage.NewRoleRepo(db))
}
func (a *app) PermissionService(ctx context.Context) userPort.PermissionService {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.permissionService == nil {
			a.permissionService = a.permissionServiceWithDB(a.db)
		}
		return a.permissionService
	}

	return a.permissionServiceWithDB(db)
}

func (a *app) permissionServiceWithDB(db *gorm.DB) userPort.PermissionService {
	return user.NewPermissionService(storage.NewPermissionRepo(db), a.surveyService)
}

func (a *app) votingServiceWithDB(db *gorm.DB, secretDB *gorm.DB) votingPort.Service {
	return voting.NewVotingService(storage.NewVotingRepo(db, secretDB))
}

func (a *app) VotingService(ctx context.Context) votingPort.Service {
	db := appCtx.GetDB(ctx)
	secretDB := appCtx.GetSecretDB(ctx)
	if db == nil || secretDB == nil {
		if a.votingService == nil {
			a.votingService = a.votingServiceWithDB(a.db, a.secretsDB)
		}
		return a.votingService
	}

	return a.votingServiceWithDB(db, secretDB)
}

func (a *app) Config(ctx context.Context) config.Config {
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
	a.codeVrfctnService = codeVerification.NewService(a.userService, storage.NewOutboxRepo(a.db), storage.NewCodeVerificationRepo(a.db))
	a.roleService = user.NewRoleService(storage.NewRoleRepo(a.db))
	a.permissionService = user.NewPermissionService(storage.NewPermissionRepo(a.db), a.surveyService)
	a.permissionService.SeedPermissions(context.Background(), generatePermissions())
	a.questionService = question.NewService(storage.NewQuestionRepo(a.db), a.surveyService)
	a.surveyService = survey.NewService(storage.NewSurveyRepo(a.db), a.permissionService)
	a.notifService = notification.NewService(storage.NewNotifRepo(a.db))
	a.votingService = voting.NewVotingService(storage.NewVotingRepo(a.db, a.secretsDB))
	return a, a.registerOutboxHandlers()
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

func (a *app) registerOutboxHandlers() error {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	common.RegisterOutboxRunner(a.codeVerificationServiceWithDB(a.db), scheduler)

	scheduler.Start()

	return nil
}
