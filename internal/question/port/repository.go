package port

import (
	"context"

	"github.com/porseOnline/pkg/adapters/storage/types"
	"gorm.io/gorm"
)

type Repo interface {
	Create(ctx context.Context, question types.Question, tx *gorm.DB) (*types.Question, error)
	Update(ctx context.Context, question types.Question, tx *gorm.DB) (*types.Question, error)
	GetNextQuestionOrder(ctx context.Context, surveyID uint) (int, error)
	Delete(ctx context.Context, id uint) (error)
	Get(ctx context.Context, id uint) (*types.Question, error)
	DeleteQuestionOptions(ctx context.Context, questionID uint, tx *gorm.DB) error
	CreateQuestionOptions(ctx context.Context, options []types.QuestionOption, questionID uint, tx *gorm.DB) ([]types.QuestionOption, error)
	GetDB(ctx context.Context) (*gorm.DB)
}