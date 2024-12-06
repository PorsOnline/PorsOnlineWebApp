package storage

import (
	"context"
	"errors"
	"time"

	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/pkg/adapters/storage/mapper"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"github.com/porseOnline/pkg/logger"
	"gorm.io/gorm"
)

type permissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) port.PermissionRepo {
	return &permissionRepo{db}
}

func (r *permissionRepo) Create(ctx context.Context, permissionDomain domain.Permission) (domain.PermissionID, error) {
	permission := mapper.PermissionDomain2Storage(permissionDomain)
	return domain.PermissionID(permission.ID), r.db.Table("permissions").WithContext(ctx).Create(permission).Error
}

func (r *permissionRepo) GetByID(ctx context.Context, permissionID domain.PermissionID) (*domain.Permission, error) {
	var permission types.Permission
	err := r.db.Debug().Table("permissions").Where("id = ?", permissionID).WithContext(ctx).First(&permission).Error

	if err != nil {
		return nil, err
	}

	if permission.ID == 0 {
		return nil, errors.New("permission not found")
	}

	return mapper.PermissionStorage2Domain(permission), nil
}

func (r *permissionRepo) Update(ctx context.Context, permissionDomain domain.Permission) error {
	var updatingPermission types.Permission
	err := r.db.Debug().Table("permissions").Where("id = ?", permissionDomain.ID).WithContext(ctx).First(&updatingPermission).Error
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Error(tx.Error.Error(), nil)
		return tx.Error
	}

	if err := tx.Model(&types.Permission{}).Where("id = ?", permissionDomain.ID).Updates(permissionDomain).Error; err != nil {
		logger.Error(err.Error(), nil)
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *permissionRepo) Delete(ctx context.Context, permissionID domain.PermissionID) error {
	return r.db.Table("permissions").Where("id = ?", permissionID).Delete(&types.Permission{}).Error
}

func (r *permissionRepo) Assign(ctx context.Context, permissionID domain.PermissionID, userID domain.UserID) error {
	var permission types.Permission
	err := r.db.Debug().Table("permissions").Where("id = ?", permissionID).WithContext(ctx).First(&permission).Error

	if err != nil {
		return err
	}

	if permission.ID == 0 {
		return errors.New("permission not found")
	}

	// find user function
	var user types.User
	err = r.db.Debug().Table("users").Where("id = ?", userID).WithContext(ctx).Preload("Role").First(&user).Error

	if err != nil {
		return err
	}

	if user.ID == 0 {
		return errors.New("user not found")
	}

	if user.Role.AccessLevel >= permission.Policy {
		//update user permissions
		tx := r.db.WithContext(ctx).Begin()
		if tx.Error != nil {
			logger.Error(tx.Error.Error(), nil)
			return tx.Error
		}

		if err := tx.Model(&user).Association("Permissions").Append(&permission); err != nil {
			logger.Error(err.Error(), nil)
			tx.Rollback()
			return err
		}

		return tx.Commit().Error
	} else {
		return errors.New("cannot assign this permission to the user")
	}
}

func (r *permissionRepo) GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Permission, error) {
	var user types.User
	err := r.db.Debug().Table("users").Where("id = ?", userID).WithContext(ctx).Preload("Permissions").First(&user, userID).Error

	if err != nil {
		return nil, err
	}

	var permissions []domain.Permission
	for _, permission := range user.Permissions {
		mappedPermission := mapper.PermissionStorage2Domain(permission)
		permissions = append(permissions, *mappedPermission)
	}
	return &permissions, nil
}

func (r *permissionRepo) Validate(ctx context.Context, userID domain.UserID, resource, scope, group string) (bool, error) {
	var user types.User
	err := r.db.Debug().Table("users").Where("id = ?", userID).WithContext(ctx).Preload("Role").Preload("Permissions").First(&user, userID).Error

	if err != nil {
		return false, err
	}

	valid := false
	for _, foundPerm := range user.Permissions {
		if foundPerm.Owner == user.ID {
			valid = true
			break
		} else if foundPerm.Resource == resource && foundPerm.Scope == scope && foundPerm.Group == group && foundPerm.CreatedAt.Add(foundPerm.Duration).Before(time.Now()) {
			valid = true
			break
		}
	}

	return valid, nil
}
