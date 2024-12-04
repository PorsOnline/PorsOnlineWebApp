package port

import (
	"context"

	"github.com/porseOnline/internal/question/domain"
)

type Service interface {
	CreateQuestion(ctx context.Context, question domain.Question) (domain.Question, error)
}