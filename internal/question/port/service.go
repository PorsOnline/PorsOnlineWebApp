package port

import (
	"context"

	"github.com/porseOnline/internal/question/domain"
)

type Service interface {
	CreateQuestion(ctx context.Context, question domain.Question) (domain.Question, error)
	DeleteQuestion(ctx context.Context, id uint) (error)
	GetQuestion(ctx context.Context, id uint) (*domain.Question, error)
}