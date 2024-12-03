package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	PasswordHash string
	NationalCode string

	BirthDate         time.Time
	City              string
	Gender            string
	SurveyLimitNumber int
	Create_at         time.Time
	Update_at         time.Time
	Balane            int
}
