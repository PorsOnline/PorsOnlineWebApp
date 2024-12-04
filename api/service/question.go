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
	if question.QuestionType == types.MultipleChoiceWithAnswer {
		if question.CorrectAnswer == "" {
			return domain.Question{}, errors.New("please choose the question answer")
		}
		for i, option := range question.QuestionOptions {
			if option.OptionText == question.CorrectAnswer {
				break
			}
			if i == len(question.QuestionOptions)-1 {
				return domain.Question{}, errors.New("please choose the question answer correctly")
			}
		}
	} else if question.QuestionType == types.ConditionalWithAnswer {
		if !(question.CorrectAnswer == "true" || question.CorrectAnswer == "false") {
			return domain.Question{}, errors.New("please choose the question answer")
		}
	}
	if question.QuestionType != types.MultipleChoice && question.QuestionType != types.MultipleChoiceWithAnswer {
		question.QuestionOptions = nil
	}
	return q.srv.CreateQuestion(ctx, *question)
}
