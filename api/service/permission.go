package service

import (
	"context"

	"github.com/porseOnline/internal/user/domain"
	userPort "github.com/porseOnline/internal/user/port"
)

type PermissionService struct {
	svc                   userPort.PermissionService
	authSecret            string
	expMin, refreshExpMin uint
}

func NewPermissionService(svc userPort.PermissionService, authSecret string, expMin, refreshExpMin uint) *PermissionService {
	return &PermissionService{svc: svc, authSecret: authSecret, expMin: expMin, refreshExpMin: refreshExpMin}
}

func (ps *PermissionService) CreatePermission(ctx context.Context, permission domain.Permission) (domain.PermissionID, error) {
	return ps.svc.CreatePermission(ctx, permission)
}

func (ps *PermissionService) GetUserPermissions(ctx context.Context, userID domain.UserID) ([]domain.Permission, error) {
	return ps.svc.GetUserPermissions(ctx, userID)
}

func (ps *PermissionService) GetPermissionByID(ctx context.Context, permissionID domain.PermissionID) (*domain.Permission, error) {
	return ps.svc.GetPermissionByID(ctx, permissionID)
}

func (ps *PermissionService) UpdatePermission(ctx context.Context, permission domain.Permission) error {
	return ps.svc.UpdatePermission(ctx, permission)
}

func (ps *PermissionService) DeletePermission(ctx context.Context, permissionID domain.PermissionID) error {
	return ps.svc.DeletePermission(ctx, permissionID)
}

func (ps *PermissionService) ValidateUserPermission(ctx context.Context, userID domain.UserID, resource, scope, group string) (bool, error) {
	return ps.svc.ValidateUserPermission(ctx, userID, resource, scope, group)
}

func (ps *PermissionService) AssignPermissionToUser(ctx context.Context, permissionID domain.PermissionID, userID domain.UserID) error {
	return ps.svc.AssignPermissionToUser(ctx, permissionID, userID)
}
