package postgres

import (
	"fmt"

	"github.com/porseOnline/pkg/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() *gorm.DB {
	host := helper.GetConfig("POSTGRES_HOST")
	user := helper.GetConfig("POSTGRES_USER")
	password := helper.GetConfig("POSTGRES_PASSWORD")
	dbname := helper.GetConfig("POSTGRES_DB_NAME")
	port := helper.GetConfig("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("Error connecting to database")
	}
	return db
}
