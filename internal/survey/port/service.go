package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/porseOnline/internal/survey/domain"
)

type Service interface {
	CreateSurvey(ctx context.Context, survey domain.Survey, userID uint) (*domain.Survey, error)
	UpdateSurvey(ctx context.Context, survey domain.Survey, id uint) (*domain.Survey, error)
	GetAllSurveys(ctx context.Context, page, pageSize int) ([]domain.Survey, error)
	GetSurveyByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Survey, error)
	GetSurveyByID(ctx context.Context, surveyID uint) (*domain.Survey, error)
	CancelSurvey(ctx context.Context, id uint) error
	DeleteSurvey(ctx context.Context, id uint) error
}
