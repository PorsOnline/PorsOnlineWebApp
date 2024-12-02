package service

import (
	"PorsOnlineWebApp/internal/survey/domain"
	surveyPort "PorsOnlineWebApp/internal/survey/port"
	"context"
)

type SurveyService struct {
	srv                   surveyPort.Service
	authSecret            string
	expMin, refreshExpMin uint
}

func NewSurveyService(srv surveyPort.Service, authSecret string, expMin, refreshExpMin uint) *SurveyService {
	return &SurveyService{srv: srv, authSecret: authSecret, expMin: expMin, refreshExpMin: refreshExpMin}
}

func (s *SurveyService) CreateSurvey(ctx context.Context, survey *domain.Survey) (*domain.Survey, error) {
	//validation
	return s.srv.CreateSurvey(ctx, *survey)
}