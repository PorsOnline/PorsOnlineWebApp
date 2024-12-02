package survey

import (
	"PorsOnlineWebApp/internal/survey/domain"
	surveyPort "PorsOnlineWebApp/internal/survey/port"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type service struct {
	repo surveyPort.Repo
}

func NewSurveyService(repo surveyPort.Repo) surveyPort.Service {
	return &service{repo: repo}
}

func (ss *service) CreateSurvey(ctx context.Context, survey domain.Survey) (*domain.Survey, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	survey.UUID = newUUID
	typeSurvey := domain.DomainToTypeMapper(survey)
	createdSurvey, err := ss.repo.Create(ctx, typeSurvey)
	if err != nil {
		return nil, err
	}
	return domain.TypeToDomainMapper(*createdSurvey), nil
}

func (ss *service) UpdateSurvey(ctx context.Context, survey domain.Survey) (*domain.Survey, error) {
	typeSurvey := domain.DomainToTypeMapper(survey)
	updatedSurvey, err := ss.repo.Update(ctx, typeSurvey)
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

func (ss *service) GetSurvey(ctx context.Context, surveyUUID uuid.UUID) (*domain.Survey, error) {
	survey, err := ss.repo.Get(ctx, surveyUUID)
	if err != nil {
		return nil, err
	}
	return domain.TypeToDomainMapper(*survey), nil
}

func (ss *service) CancelSurvey(ctx context.Context, surveyUUID uuid.UUID) error {
	return ss.repo.Cancel(ctx, surveyUUID)
}

func (ss *service) DeleteSurvey(ctx context.Context, surveyUUID uuid.UUID) error {
	return ss.repo.Cancel(ctx, surveyUUID)
}
