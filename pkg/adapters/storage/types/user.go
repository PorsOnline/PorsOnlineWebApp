package types

import (
	"time"
)

type User struct {
	ID                string
	FirstName         string
	Lastname          string
	NationalCode      string
	Email             string
	Password          string
	BirthDate         time.Time
	City              string
	Sex               string
	SurveyLimitNumber int
	Create_at         time.Time
	Update_at         time.Time
	Balane            int
}
