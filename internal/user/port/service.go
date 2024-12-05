package port

import (
	"context"

	"github.com/porseOnline/internal/user/domain"
)

type Service interface {
	CreateUser(ctx context.Context, user domain.User) (domain.UserID, error)
	GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email domain.Email) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
}
