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
	DeleteByID(ctx context.Context, userID domain.UserID) error
}

type RoleService interface {
	CreateRole(ctx context.Context, role domain.Role) (domain.RoleID, error)
	GetRole(ctx context.Context, roleID domain.RoleID) (*domain.Role, error)
	UpdateRole(ctx context.Context, role domain.Role) error
	DeleteRole(ctx context.Context, roleID domain.RoleID) error
	AssignRoleToUser(ctx context.Context, roleID domain.RoleID, userID domain.UserID) error
}

type PermissionService interface {
	CreatePermission(ctx context.Context, permission domain.Permission) (domain.PermissionID, error)
	GetUserPermissions(ctx context.Context, userID domain.UserID) ([]domain.Permission, error)
	GetPermissionByID(ctx context.Context, permissionID domain.PermissionID) (*domain.Permission, error)
	UpdatePermission(ctx context.Context, permission domain.Permission) error
	DeletePermission(ctx context.Context, permissionID domain.PermissionID) error
	ValidateUserPermission(ctx context.Context, userID domain.UserID, resource, scope, group string) (bool, error)
	AssignPermissionToUser(ctx context.Context,permissionDetails []domain.PermissionDetails) error
	SeedPermissions(ctx context.Context, permissions []domain.Permission) error 
}
