package service

import (
	"context"
	"errors"

	"github.com/porseOnline/internal/question/domain"
	questionPort "github.com/porseOnline/internal/question/port"
	"github.com/porseOnline/pkg/adapters/storage/types"
)

type QuestionService struct {
	srv                   questionPort.Service
	authSecret            string
	expMin, refreshExpMin uint
}

func NewQuestionService(srv questionPort.Service, authSecret string, expMin, refreshExpMin uint) *QuestionService {
	return &QuestionService{srv: srv, authSecret: authSecret, expMin: expMin, refreshExpMin: refreshExpMin}
}

func (q *QuestionService) CreateQuestion(ctx context.Context, question *domain.Question) (domain.Question, error) {
	err := validateQuestionType(*question)
	if err != nil {
		return domain.Question{}, err
	}
	if question.QuestionType == types.Descriptive {
		question.QuestionOptions = nil
	}
	return q.srv.CreateQuestion(ctx, *question)
}

func (q *QuestionService) DeleteQuestion(ctx context.Context, id uint, surveyID uint) (error) {
	return q.srv.DeleteQuestion(ctx, id, surveyID)
}

func (q *QuestionService) GetNextQuestion(ctx context.Context, questionStep domain.UserQuestionStep, userID uint, surveyID uint) (*domain.Question, error) {
	return q.srv.GetNextQuestion(ctx, questionStep, userID, surveyID)
}

func (q *QuestionService) UpdateQuestion(ctx context.Context, question *domain.Question) (domain.Question, error) {
	err := validateQuestionType(*question)
	if err != nil {
		return domain.Question{}, err
	}
	if question.QuestionType == types.Descriptive {
		question.QuestionOptions = nil
	}
	return q.srv.CreateQuestion(ctx, *question)
}

func validateQuestionType(question domain.Question) error {
	var countCorrectAnswers int
	if question.QuestionType == types.MultipleChoiceWithAnswer{
		for _, option := range question.QuestionOptions {
			if option.IsCorrect {
				countCorrectAnswers += 1
			}
		}
		if countCorrectAnswers != 1 {
			return errors.New("choose correct answer correctly")
		}
	}
	return nil
}