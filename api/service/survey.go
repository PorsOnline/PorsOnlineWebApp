package service

import (
	"github.com/porseOnline/internal/survey/domain"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"context"
)

type SurveyService struct {
	srv                   surveyPort.Service
	authSecret            string
	expMin, refreshExpMin uint
}

func NewService(srv surveyPort.Service, authSecret string, expMin, refreshExpMin uint) *SurveyService {
	return &SurveyService{srv: srv, authSecret: authSecret, expMin: expMin, refreshExpMin: refreshExpMin}
}

func (s *SurveyService) CreateSurvey(ctx context.Context, survey *domain.Survey, userID uint) (*domain.Survey, error) {
	//validation
	return s.srv.CreateSurvey(ctx, *survey, userID)
}

func (s *SurveyService) GetSurvey(ctx context.Context, id uint) (*domain.Survey, error) {
	return s.srv.GetSurveyByID(ctx, id)
}

func (s *SurveyService) UpdateSurvey(ctx context.Context, survey *domain.Survey, id uint) (*domain.Survey, error) {
	//validation
	return s.srv.UpdateSurvey(ctx, *survey, id)
}

func (s *SurveyService) CancelSurvey(ctx context.Context, id uint) error {
	return s.srv.CancelSurvey(ctx, id)
}

func (s *SurveyService) DeleteSurvey(ctx context.Context, id uint) error {
	return s.srv.DeleteSurvey(ctx, id)
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