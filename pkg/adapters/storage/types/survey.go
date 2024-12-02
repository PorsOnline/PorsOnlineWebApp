package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Survey struct {
	gorm.Model
	UUID uuid.UUID
	Title string
	StartAt time.Time
	ExpireAt time.Time
	IsSequential bool
	IsActive bool
	AllowsBackNavigation bool
	MaxAttempts uint
	DurationMinutes uint
}