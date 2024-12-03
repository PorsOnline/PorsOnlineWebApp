package domain

import (
	"time"

	"github.com/porseOnline/pkg/adapters/storage/types"

	"github.com/google/uuid"
)

type Survey struct {
	ID uint
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

func TypeToDomainMapper(survey types.Survey) *Survey {
	return &Survey{
		survey.ID,
		survey.UUID,
		survey.Title,
		survey.StartAt,
		survey.ExpireAt,
		survey.IsSequential,
		survey.IsActive,
		survey.AllowsBackNavigation,
		survey.MaxAttempts,
		survey.DurationMinutes,
	}
}

func DomainToTypeMapper(survey Survey) types.Survey {
	return types.Survey {
		UUID: survey.UUID,
		Title: survey.Title,
		StartAt: survey.StartAt,
		ExpireAt: survey.ExpireAt,
		IsSequential: survey.IsSequential,
		IsActive: survey.IsActive,
		AllowsBackNavigation: survey.AllowsBackNavigation,
		MaxAttempts: survey.MaxAttempts,
		DurationMinutes: survey.DurationMinutes,
	}
}