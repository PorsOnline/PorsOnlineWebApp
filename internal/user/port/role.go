package port

import (
	"context"

	"github.com/porseOnline/internal/user/domain"
)

type RoleRepo interface {
	Create(ctx context.Context, role domain.Role) (domain.RoleID, error)
	Get(ctx context.Context, roleID domain.RoleID) (*domain.Role, error)
	Update(ctx context.Context, role domain.Role) error
	Delete(ctx context.Context, roleID domain.RoleID) error
	Assign(ctx context.Context, roleID domain.RoleID, userID domain.UserID) error
}
