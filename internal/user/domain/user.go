package domain

import (
	"errors"
	"time"
)

type (
	UserID uint
	Phone  string
	Email  string
)

func (p Phone) IsValid() bool {
	// todo regex
	return true
}

type User struct {
	ID                UserID
	FirstName         string
	LastName          string
	Phone             Phone
	Email             Email
	PasswordHash      string
	NationalCode      string
	BirthDate         time.Time
	City              string
	Gender            bool
	SurveyLimitNumber int
	CreatedAt         time.Time
	DeletedAt         time.Time
	UpdatedAt         time.Time
	Balance           int
}

func (u *User) Validate() error {
	if !u.Phone.IsValid() {
		return errors.New("phone is not valid")
	}
	return nil
}
