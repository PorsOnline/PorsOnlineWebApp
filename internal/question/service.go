package question

import (
	"context"
	"errors"

	"github.com/porseOnline/internal/question/domain"
	questionPort "github.com/porseOnline/internal/question/port"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"gorm.io/gorm"
)

type questionService struct {
	questionRepo questionPort.Repo
	surveyService surveyPort.Service
}

func NewService(repo questionPort.Repo, surveyService surveyPort.Service) questionPort.Service {
	return &questionService{questionRepo: repo, surveyService: surveyService}
}

func (qs *questionService) CreateQuestion(ctx context.Context, question domain.Question) (domain.Question, error) {
	survey, err := qs.surveyService.GetSurvey(ctx, question.SurveyUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Question{}, errors.New("survey not found")
		}
		return domain.Question{}, err
	}
	questionType := domain.DomainToTypeMapper(question, survey.ID)
	if !question.IsDependency {
		order, err := qs.questionRepo.GetNextQuestionOrder(ctx, questionType.SurveyID)
		if err != nil {
			return domain.Question{}, err
		}
		questionType.Order = order
	}
	createdQuestion, err := qs.questionRepo.Create(ctx, questionType)
	if err != nil {
		return domain.Question{}, err
	}
	return *domain.TypeToDomainMapper(*createdQuestion, survey.UUID), nil
}
