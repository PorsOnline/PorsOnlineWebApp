package postgres

import (
	"fmt"

	"github.com/porseOnline/pkg/adapters/storage/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnOptions struct {
	User   string
	Pass   string
	Host   string
	Port   uint
	DBName string
	Schema string
}

func (o DBConnOptions) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		o.Host, o.Port, o.User, o.Pass, o.DBName, o.Schema)
}

func NewPsqlGormConnection(opt DBConnOptions) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(opt.PostgresDSN()), &gorm.Config{
		Logger: logger.Default,
	})
}
func GormMigrations(db *gorm.DB) {
	db.AutoMigrate(
		&types.Permission{},
		&types.UserPermission{},
		&types.Notification{},
		&types.User{},
		&types.CodeVerification{},
		&types.Vote{},
		&types.Survey{},
		&types.SurveyCity{},
		&types.Question{},
		&types.QuestionOption{},
	)
}

func GormSecretsMigration(db *gorm.DB) {
	db.AutoMigrate(
		&types.Secrets{},
	)
}
