package port

import (
	"context"

	"github.com/porseOnline/internal/question/domain"
)

type Service interface {
	CreateQuestion(ctx context.Context, question domain.Question) (domain.Question, error)
	UpdateQuestion(ctx context.Context, question domain.Question) (domain.Question, error)
	DeleteQuestion(ctx context.Context, id uint) (error)
	GetQuestionByID(ctx context.Context, id uint) (*domain.Question, error)
	GetNextQuestion(ctx context.Context, user domain.UserQuestionStep, userID uint) (*domain.Question, error)
}