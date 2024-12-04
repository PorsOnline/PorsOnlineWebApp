package port

import (
	"github.com/porseOnline/internal/survey/domain"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateSurvey(ctx context.Context, survey domain.Survey) (*domain.Survey, error)
	UpdateSurvey(ctx context.Context, survey domain.Survey) (*domain.Survey, error)
	GetAllSurveys(ctx context.Context, page, pageSize int) ([]domain.Survey, error)
	GetSurveyByUUID(ctx context.Context, surveyUUID uuid.UUID) (*domain.Survey, error)
	GetSurveyByID(ctx context.Context, surveyID uint) (*domain.Survey, error)
	CancelSurvey(ctx context.Context, surveyUUID uuid.UUID) error
	DeleteSurvey(ctx context.Context, surveyUUID uuid.UUID) error
}
