package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Survey struct {
	gorm.Model
	UUID                 uuid.UUID
	Title                string
	StartAt              time.Time
	ExpireAt             time.Time
	IsSequential         bool
	IsActive             bool
	AllowsBackNavigation bool
	MaxAttempts          uint
	DurationMinutes      uint
	MinAge               uint
	MaxAge               uint
	Gender               bool
	TargetCities         []SurveyCity
}

type GenderEnum int

const (
	Male GenderEnum = iota
	Female
	Unknown
)

type SurveyCity struct {
	ID       uint `gorm:"primarykey"`
	SurveyID uint
	Name     string
}
