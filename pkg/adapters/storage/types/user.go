package types

import (
	"time"

	"github.com/porseOnline/internal/user/domain"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName         string    `gorm:"column:first_name"`
	LastName          string    `gorm:"column:last_name"`
	Phone             string    `gorm:"column:phone;unique;not null"`
	Email             string    `gorm:"column:email;unique;not null"`
	PasswordHash      string    `gorm:"column:password_hash;not null"`
	NationalCode      string    `gorm:"column:national_code;unique;not null"`
	BirthDate         time.Time `gorm:"column:birth_date"`
	City              string    `gorm:"column:city"`
	Gender            bool      `gorm:"column:gender;not null"`
	SurveyLimitNumber int       `gorm:"column:survey_limit_number;default:100"`
	Balance           int       `gorm:"column:balance;default:100"`
	RoleID            *uint
	Role              *Role `gorm:"foreignkey:RoleID"`
	UserPermissions   []UserPermission
}

type UserPermission struct {
	gorm.Model
	UserID       uint
	User         *User `gorm:"foreignkey:UserID"`
	PermissionID domain.PermissionID
	Permission   Permission `gorm:"foreignkey:PermissionID"`
	SurveyID     *uint
	Survey       *Survey `gorm:"foreignkey:SurveyID"`
	Duration     time.Duration
}
