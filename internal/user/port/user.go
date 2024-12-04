package port

import (
	"context"

	"github.com/porseOnline/internal/user/domain"
)

type Repo interface {
	Create(ctx context.Context, user domain.User) (domain.UserID, error)
	GetByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
	GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error)
}
