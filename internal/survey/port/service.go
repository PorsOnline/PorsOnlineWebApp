package port

import (
	"PorsOnlineWebApp/internal/survey/domain"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateSurvey(ctx context.Context, survey domain.Survey) (*domain.Survey, error)
	UpdateSurvey(ctx context.Context, survey domain.Survey) (*domain.Survey, error)
	GetAllSurveys(ctx context.Context, page, pageSize int) ([]domain.Survey, error)
	GetSurvey(ctx context.Context, surveyUUID uuid.UUID) (*domain.Survey, error)
	CancelSurvey(ctx context.Context, surveyUUID uuid.UUID) error
	DeleteSurvey(ctx context.Context, surveyUUID uuid.UUID) error
}
