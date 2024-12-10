package service

import (
	"context"

	"github.com/porseOnline/internal/user/domain"
	userPort "github.com/porseOnline/internal/user/port"
)

type RoleService struct {
	svc                   userPort.RoleService
	authSecret            string
	expMin, refreshExpMin uint
}

func NewRoleService(svc userPort.RoleService, authSecret string, expMin, refreshExpMin uint) *RoleService {
	return &RoleService{svc: svc, authSecret: authSecret, expMin: expMin, refreshExpMin: refreshExpMin}
}

func (rs *RoleService) CreateRole(ctx context.Context, role domain.Role) (domain.RoleID, error) {
	return rs.svc.CreateRole(ctx, role)
}

func (rs *RoleService) GetRole(ctx context.Context, roleID domain.RoleID) (*domain.Role, error) {
	return rs.svc.GetRole(ctx, roleID)
}

func (rs *RoleService) UpdateRole(ctx context.Context, role domain.Role) error {
	return rs.svc.UpdateRole(ctx, role)
}

func (rs *RoleService) DeleteRole(ctx context.Context, roleID domain.RoleID) error {
	return rs.svc.DeleteRole(ctx, roleID)
}

func (rs *RoleService) AssignRoleToUser(ctx context.Context, roleID domain.RoleID, userID domain.UserID) error {
	return rs.svc.AssignRoleToUser(ctx, roleID, userID)
}
