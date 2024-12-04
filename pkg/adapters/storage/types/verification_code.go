package types

import (
	"time"

	"gorm.io/gorm"
)

type CodeVerification struct {
	gorm.Model
	UserID uint      `gorm:"column:userId"`
	Code   string    `gorm:"column:code"`
	ExpMin time.Time `gorm:"column:expMin"`
}
