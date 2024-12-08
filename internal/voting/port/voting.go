package port

import (
	"context"

	"github.com/porseOnline/internal/voting/domain"
)

type Repo interface {
	Vote(ctx context.Context, answer *domain.Vote) error
	GetLastResponse(ctx context.Context, userID uint, serveyID uint) (domain.Vote, error)
}