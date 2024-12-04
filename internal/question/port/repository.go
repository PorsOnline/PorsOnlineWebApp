package port

import (
	"context"

	"github.com/porseOnline/pkg/adapters/storage/types"
)

type Repo interface {
	Create(ctx context.Context, question types.Question) (*types.Question, error)
	GetNextQuestionOrder(ctx context.Context, surveyID uint) (int, error)
	Delete(ctx context.Context, id uint) (error)
	Get(ctx context.Context, id uint) (*types.Question, error)
}