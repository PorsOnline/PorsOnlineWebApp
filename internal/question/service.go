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
	survey, err := qs.surveyService.GetSurveyByUUID(ctx, question.SurveyUUID)
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

func (qs *questionService) DeleteQuestion(ctx context.Context, id uint) (error) {
	err := qs.questionRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (qs *questionService) GetQuestion(ctx context.Context, id uint) (*domain.Question, error) {
	question, err := qs.questionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	survey, err := qs.surveyService.GetSurveyByID(ctx, question.ID)
	if err != nil {
		return nil, err
	}
	return domain.TypeToDomainMapper(*question, survey.UUID), nil
}
