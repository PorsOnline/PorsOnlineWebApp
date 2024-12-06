package types

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Owner    uint   `gorm:"column:owner"`
	Group    string `gorm:"column:group"`
	Resource string `gorm:"column:resource"`
	Scope    string `gorm:"column:scope"`
	Policy   uint8  `gorm:"column:policy"`
	Users    []User `gorm:"many2many:user_permissions"`
}
