package types

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Owner    uint          `gorm:"column:owner"`
	Group    string        `gorm:"column:group"`
	Resource string        `gorm:"column:resource"`
	Scope    string        `gorm:"column:scope"`
	Policy   uint8         `gorm:"column:policy"`
	Duration time.Duration `gorm:"column:duration"`
	Users    []User        `gorm:"many2many:user_permissions"`
}
