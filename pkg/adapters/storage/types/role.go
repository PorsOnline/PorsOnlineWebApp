package types

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"column:name"`
	AccessLevel uint8  `gorm:"column:access_level"`
}
