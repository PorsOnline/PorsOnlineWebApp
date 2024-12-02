//	func NewConnection() *gorm.DB {
//		host := helper.GetConfig("POSTGRES_HOST")
//		user := helper.GetConfig("POSTGRES_USER")
//		password := helper.GetConfig("POSTGRES_PASSWORD")
//		dbname := helper.GetConfig("POSTGRES_DB_NAME")
//		port := helper.GetConfig("POSTGRES_PORT")
//		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", host, user, password, dbname, port)
//		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//		if err != nil {
//			fmt.Println(err)
//			panic("Error connecting to database")
//		}
//		return db

package postgres

import (
	"fmt"

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
