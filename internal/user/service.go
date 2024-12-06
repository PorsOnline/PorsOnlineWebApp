package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/pkg/logger"
)

var (
	ErrUserOnCreate           = errors.New("error on creating new user")
	ErrUserCreationValidation = errors.New("validation failed")
	ErrUserNotFound           = errors.New("user not found")
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, user domain.User) (domain.UserID, error) {
	if err := user.Validate(); err != nil {
		return 0, fmt.Errorf("%w %w", ErrUserCreationValidation, err)
	}

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		log.Println("error on creating new user : ", err.Error())
		return 0, ErrUserOnCreate
	}

	return userID, nil
}

func (s *service) GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == 0 {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == 0 {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, user domain.User) error {
	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		logger.Error("error in update user", nil)
		return err
	}
	return nil
}

func (s *service) DeleteByID(ctx context.Context, userID domain.UserID) error {
	err := s.repo.DeleteByID(ctx, userID)
	if err != nil {
		logger.Error("can not delete user", nil)
		return err
	}
	return nil
}

// ----------------- Role related services
type roleService struct {
	repo port.RoleRepo
}

func NewRoleService(roleRepo port.RoleRepo) port.RoleService {
	return &roleService{
		repo: roleRepo,
	}
}

func (rs *roleService) CreateRole(ctx context.Context, role domain.Role) (domain.RoleID, error) {
	roleID, err := rs.repo.Create(ctx, role)
	if err != nil {
		logger.Error("error on creating new role", nil)
		return 0, err
	}
	logger.Info("successful create role", nil)
	return roleID, nil
}

func (rs *roleService) GetRole(ctx context.Context, roleID domain.RoleID) (*domain.Role, error) {
	role, err := rs.repo.Get(ctx, roleID)
	if err != nil {
		logger.Error("role not found", nil)
		return nil, err
	}
	if role == nil || role.ID == 0 {
		logger.Error("role not found", nil)
		return nil, errors.New("role no found")
	}
	logger.Info("successful get role", nil)
	return role, nil
}

func (rs *roleService) UpdateRole(ctx context.Context, role domain.Role) error {
	err := rs.repo.Update(ctx, role)
	if err != nil {
		logger.Error("error in updating role", nil)
		return err
	}
	logger.Info("successful update role", nil)
	return nil
}

func (rs *roleService) DeleteRole(ctx context.Context, roleID domain.RoleID) error {
	err := rs.repo.Delete(ctx, roleID)
	if err != nil {
		logger.Error("error in deleting role", nil)
		return err
	}
	logger.Info("successful delete role", nil)
	return nil
}

func (rs *roleService) AssignRoleToUser(ctx context.Context, roleID domain.RoleID, userID domain.UserID) error {
	err := rs.repo.Assign(ctx, roleID, userID)
	if err != nil {
		logger.Error("error in assigning role to user", nil)
		return err
	}
	logger.Info("successful assign role to user", nil)
	return nil
}

// ----------------- Permission related services
type permissionService struct {
	repo port.PermissionRepo
}

func NewPermissionService(permissionRepo port.PermissionRepo) port.PermissionService {
	return &permissionService{
		repo: permissionRepo,
	}
}

func (ps *permissionService) CreatePermission(ctx context.Context, permission domain.Permission) (domain.PermissionID, error) {
	permissionID, err := ps.repo.Create(ctx, permission)
	if err != nil {
		logger.Error("error on creating new permission", nil)
		return 0, err
	}
	logger.Info("successful create permission", nil)
	return permissionID, nil
}

func (ps *permissionService) GetPermissionByID(ctx context.Context, permissionID domain.PermissionID) (*domain.Permission, error) {
	permission, err := ps.repo.GetByID(ctx, permissionID)
	if err != nil {
		logger.Error("permission not found", nil)
		return nil, err
	}
	if permission == nil || permission.ID == 0 {
		logger.Error("permission not found", nil)
		return nil, errors.New("permission no found")
	}
	logger.Info("successful get permission", nil)
	return permission, nil
}

func (ps *permissionService) UpdatePermission(ctx context.Context, permission domain.Permission) error {
	err := ps.repo.Update(ctx, permission)
	if err != nil {
		logger.Error("error in updating permission", nil)
		return err
	}
	logger.Info("successful update permission", nil)
	return nil
}

func (ps *permissionService) DeletePermission(ctx context.Context, permissionID domain.PermissionID) error {
	err := ps.repo.Delete(ctx, permissionID)
	if err != nil {
		logger.Error("error in deleting permission", nil)
		return err
	}
	logger.Info("successful delete permission", nil)
	return nil
}

func (ps *permissionService) AssignPermissionToUser(ctx context.Context, permissionID domain.PermissionID, userID domain.UserID) error {
	err := ps.repo.Assign(ctx, permissionID, userID)
	if err != nil {
		logger.Error("error in assigning permission to user", nil)
		return err
	}
	logger.Info("successful assign permission to user", nil)
	return nil
}

func (ps *permissionService) GetUserPermissions(ctx context.Context, userID domain.UserID) ([]domain.Permission, error) {
	permissions, err := ps.repo.GetAll(ctx, userID)
	if err != nil {
		logger.Error("error in getting user permissions", nil)
		return nil, err
	}
	logger.Info("successful list user permissions", nil)
	return *permissions, nil
}

func (ps *permissionService) ValidateUserPermission(ctx context.Context, userID domain.UserID, resource, scope, group string) (bool, error) {
	valid, err := ps.repo.Validate(ctx, userID, resource, scope, group)
	if err != nil {
		logger.Error("error in validating user access", nil)
		return valid, err
	}
	logger.Info("successful validation on user access", nil)
	return true, nil
}
