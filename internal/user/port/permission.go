package port

import (
	"context"

	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/pkg/adapters/storage/types"
)

type PermissionRepo interface {
	Create(ctx context.Context, permission domain.Permission) (domain.PermissionID, error)
	GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Permission, error)
	GetByID(ctx context.Context, permissionID domain.PermissionID) (*domain.Permission, error)
	Update(ctx context.Context, permission domain.Permission) error
	Delete(ctx context.Context, permissionID domain.PermissionID) error
	Validate(ctx context.Context, userID domain.UserID, resource, scope, group string) (bool, error)
	Assign(ctx context.Context, userPermission types.UserPermission) error
}
