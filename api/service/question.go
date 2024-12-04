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
	if question.QuestionType != types.MultipleChoice && question.QuestionType != types.MultipleChoiceWithAnswer {
		question.QuestionOptions = nil
	}
	return q.srv.CreateQuestion(ctx, *question)
}

func (q *QuestionService) DeleteQuestion(ctx context.Context, id uint) (error) {
	return q.srv.DeleteQuestion(ctx, id)
}

func (q *QuestionService) UpdateQuestion(ctx context.Context, question *domain.Question) (domain.Question, error) {
	err := validateQuestionType(*question)
	if err != nil {
		return domain.Question{}, err
	}
	if question.QuestionType != types.MultipleChoice && question.QuestionType != types.MultipleChoiceWithAnswer {
		question.QuestionOptions = nil
	}
	return q.srv.CreateQuestion(ctx, *question)
}

func validateQuestionType(question domain.Question) error {
	if question.QuestionType == types.MultipleChoiceWithAnswer {
		if question.CorrectAnswer == "" {
			return errors.New("please choose the question answer")
		}
		for i, option := range question.QuestionOptions {
			if option.OptionText == question.CorrectAnswer {
				break
			}
			if i == len(question.QuestionOptions)-1 {
				return errors.New("please choose the question answer correctly")
			}
		}
	} else if question.QuestionType == types.ConditionalWithAnswer {
		if !(question.CorrectAnswer == "true" || question.CorrectAnswer == "false") {
			return errors.New("please choose the question answer")
		}
	}
	return nil
}