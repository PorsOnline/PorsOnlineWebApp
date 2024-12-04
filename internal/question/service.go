package question

import (
	"context"

	"github.com/porseOnline/internal/question/domain"
	questionPort "github.com/porseOnline/internal/question/port"
)

type questionService struct {
	questionRepo questionPort.Repo
}

func NewQuestionService(repo questionPort.Repo) questionPort.Service {
	return &questionService{questionRepo: repo}
}

func (qs *questionService) CreateQuestion(ctx context.Context, question domain.Question) (domain.Question, error) {
	questionType := domain.DomainToTypeMapper(question)
	if !*question.IsDependency {
		order, err := qs.questionRepo.GetNextQuestionOrder(ctx, question.SurveyID)
		if err != nil {
			return domain.Question{}, err
		}
		questionType.Order = order
	}
	createdQuestion, err := qs.questionRepo.Create(ctx, questionType)
	if err != nil {
		return domain.Question{}, err
	}
	return *domain.TypeToDomainMapper(*createdQuestion), nil
}
