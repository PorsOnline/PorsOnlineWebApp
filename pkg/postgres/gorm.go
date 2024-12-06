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
		Logger: logger.Discard,
	})
}
func GormMigrations(db *gorm.DB) {
	db.AutoMigrate(
		&types.Notification{},
		&types.User{},
		&types.CodeVerification{},
	)
}

func GormSecretsMigration(db *gorm.DB) {
	db.AutoMigrate(
		&types.Secrets{},
	)
}
