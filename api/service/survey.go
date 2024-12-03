package service

import (
	"github.com/porseOnline/internal/survey/domain"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"context"

	"github.com/google/uuid"
)

type SurveyService struct {
	srv                   surveyPort.Service
	authSecret            string
	expMin, refreshExpMin uint
}

func NewService(srv surveyPort.Service, authSecret string, expMin, refreshExpMin uint) *SurveyService {
	return &SurveyService{srv: srv, authSecret: authSecret, expMin: expMin, refreshExpMin: refreshExpMin}
}

func (s *SurveyService) CreateSurvey(ctx context.Context, survey *domain.Survey) (*domain.Survey, error) {
	//validation
	return s.srv.CreateSurvey(ctx, *survey)
}

func (s *SurveyService) GetSurvey(ctx context.Context, uuid uuid.UUID) (*domain.Survey, error) {
	return s.srv.GetSurvey(ctx, uuid)
}

func (s *SurveyService) UpdateSurvey(ctx context.Context, survey *domain.Survey) (*domain.Survey, error) {
	//validation
	return s.srv.UpdateSurvey(ctx, *survey)
}

func (s *SurveyService) CancelSurvey(ctx context.Context, uuid uuid.UUID) error {
	return s.srv.CancelSurvey(ctx, uuid)
}

func (s *SurveyService) DeleteSurvey(ctx context.Context, uuid uuid.UUID) error {
	return s.srv.DeleteSurvey(ctx, uuid)
}

func (s *SurveyService) GetAllSurveys(ctx context.Context, page, pageSize int) ([]domain.Survey, error) {
	if page == 0 {
		page = 1
	} 
	if pageSize == 0 {
		pageSize = 10
	}
	return s.srv.GetAllSurveys(ctx, page, pageSize)
}