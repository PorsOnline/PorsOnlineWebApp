package port

import (
	"github.com/porseOnline/pkg/adapters/storage/types"
	"context"

	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, survey types.Survey, cities []string) (*types.Survey, error)
	Update(ctx context.Context, survey types.Survey, cities []string) (*types.Survey, error)
	GetAll(ctx context.Context, page, pageSize int) ([]types.Survey, error)
	GetByUUID(ctx context.Context, surveyUUID uuid.UUID) (*types.Survey, error)
	GetByID(ctx context.Context, id uint) (*types.Survey, error)
	Cancel(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}