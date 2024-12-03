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
	MinAge               uint
	MaxAge               uint
	Gender               types.GenderEnum
	TargetCities         []string
}

func TypeToDomainMapper(survey types.Survey) *Survey {
	var cities []string
	for _, city := range survey.TargetCities {
		cities = append(cities, city.Name)
	}
	return &Survey{
		ID: survey.ID,
		UUID: survey.UUID,
		Title: survey.Title,
		StartAt: survey.StartAt,
		ExpireAt: survey.ExpireAt,
		IsSequential: survey.IsSequential,
		IsActive: survey.IsActive,
		AllowsBackNavigation: survey.AllowsBackNavigation,
		MaxAttempts: survey.MaxAttempts,
		DurationMinutes: survey.DurationMinutes,
		MinAge: survey.MinAge,
		MaxAge: survey.MaxAge,
		Gender: survey.Gender,
		TargetCities: cities,
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
		MinAge: survey.MinAge,
		MaxAge: survey.MaxAge,
		Gender: survey.Gender,
	}
}