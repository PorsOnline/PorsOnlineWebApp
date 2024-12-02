package port

import (
	"PorsOnlineWebApp/pkg/adapters/storage/types"
	"context"

	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, survey types.Survey) (*types.Survey, error)
	Update(ctx context.Context, survey types.Survey) (*types.Survey, error)
	GetAll(ctx context.Context, page, pageSize int) ([]types.Survey, error)
	Get(ctx context.Context, surveyUUID uuid.UUID) (*types.Survey, error)
	Cancel(ctx context.Context, surveyUUID uuid.UUID) error
	Delete(ctx context.Context, surveyUUID uuid.UUID) error
}