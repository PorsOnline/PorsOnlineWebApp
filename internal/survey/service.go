package survey

import (
	"github.com/porseOnline/internal/survey/domain"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type service struct {
	repo surveyPort.Repo
}

func NewService(repo surveyPort.Repo) surveyPort.Service {
	return &service{repo: repo}
}

func (ss *service) CreateSurvey(ctx context.Context, survey domain.Survey) (*domain.Survey, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	survey.UUID = newUUID
	typeSurvey := domain.DomainToTypeMapper(survey)
	createdSurvey, err := ss.repo.Create(ctx, typeSurvey, survey.TargetCities)
	if err != nil {
		return nil, err
	}
	return domain.TypeToDomainMapper(*createdSurvey), nil
}

func (ss *service) UpdateSurvey(ctx context.Context, survey domain.Survey) (*domain.Survey, error) {
	typeSurvey := domain.DomainToTypeMapper(survey)
	updatedSurvey, err := ss.repo.Update(ctx, typeSurvey, survey.TargetCities)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("survey not found")
		}
		return nil, err
	}
	return domain.TypeToDomainMapper(*updatedSurvey), nil
}

func (ss *service) GetAllSurveys(ctx context.Context, page, pageSize int) ([]domain.Survey, error) {
	surveys, err := ss.repo.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	var domainSurveys []domain.Survey
	for _, survey := range surveys {
		domainSurveys = append(domainSurveys, *domain.TypeToDomainMapper(survey))
	}
	return domainSurveys, nil
}

func (ss *service) GetSurveyByUUID(ctx context.Context, surveyUUID uuid.UUID) (*domain.Survey, error) {
	survey, err := ss.repo.GetByUUID(ctx, surveyUUID)
	if err != nil {
		return nil, err
	}
	return domain.TypeToDomainMapper(*survey), nil
}

func (ss *service) CancelSurvey(ctx context.Context, surveyUUID uuid.UUID) error {
	return ss.repo.Cancel(ctx, surveyUUID)
}

func (ss *service) DeleteSurvey(ctx context.Context, surveyUUID uuid.UUID) error {
	return ss.repo.Delete(ctx, surveyUUID)
}

func (ss *service) GetSurveyByID(ctx context.Context, surveyID uint) (*domain.Survey, error) {
	survey, err := ss.repo.GetByID(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	return domain.TypeToDomainMapper(*survey), nil
}
