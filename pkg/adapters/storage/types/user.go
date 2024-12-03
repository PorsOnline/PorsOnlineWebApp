package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName         string    `gorm:"column:first_name"`
	LastName          string    `gorm:"column:last_name"`
	Phone             string    `gorm:"column:phone"`
	Email             string    `gorm:"column:email"`
	PasswordHash      string    `gorm:"column:password_hash"`
	NationalCode      string    `gorm:"column:national_code"`
	BirthDate         time.Time `gorm:"column:birth_date"`
	City              string    `gorm:"column:city"`
	Gender            string    `gorm:"column:gender"`
	SurveyLimitNumber int       `gorm:"column:survey_limit_number"`
	Balance           int       `gorm:"column:balance"`
}
