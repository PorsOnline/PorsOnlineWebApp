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

func (r *permissionRepo) Assign(ctx context.Context, userPermission types.UserPermission) error {
	var permission types.Permission
	err := r.db.Debug().Table("permissions").Where("id = ?", userPermission.PermissionID).WithContext(ctx).First(&permission).Error

	if err != nil {
		return err
	}

	if permission.ID == 0 {
		return errors.New("permission not found")
	}

	// find user function
	var user types.User
	err = r.db.Debug().Table("users").Where("id = ?", userPermission.UserID).WithContext(ctx).Preload("Role").First(&user).Error

	if err != nil {
		return err
	}

	if user.ID == 0 {
		return errors.New("user not found")
	}

	if user.Role.AccessLevel >= permission.Policy {
		//update user permissions
		if err := r.db.Model(&types.UserPermission{}).Create(&userPermission).Error; err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("cannot assign this permission to the user")
	}
}

func (r *permissionRepo) GetAll(ctx context.Context, userID domain.UserID) (*[]domain.Permission, error) {
	var user types.User
	err := r.db.Debug().Table("users").Where("id = ?", userID).WithContext(ctx).Preload("UserPermissions").First(&user, userID).Error

	if err != nil {
		return nil, err
	}

	var permissions []domain.Permission
	for _, userPermission := range user.UserPermissions {
		permission, err := r.GetByID(ctx, domain.PermissionID(userPermission.PermissionID))
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, *permission)
	}
	return &permissions, nil
}

func (r *permissionRepo) Validate(ctx context.Context, userID domain.UserID, resource, scope, group string) (bool, error) {
	var userPermissionDetails types.UserPermission
	err := r.db.Table("users u").
		Joins("left join user_permissions up on u.id = up.user_id").
		Joins("left join permissions p on p.id = up.permission_id").
		Select("up.id", "up.duration", "up.created_at").
		Where("u.id = ? and (? like replace(p.resource, ':id', '%') or ? like replace(p.resource, ':uuid', '%')) and p.scope = ?", userID, resource, resource, scope).
		First(&userPermissionDetails).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if userPermissionDetails.ID > 0 {
		if userPermissionDetails.CreatedAt.Add(userPermissionDetails.Duration).Before(time.Now()) {
			return true, nil
		}
	}
	return false, nil

	// for _, userPermission := range user.UserPermissions {
	// 	if userPermission.Permission.Owner == user.ID {
	// 		valid = true
	// 		break
	// 	} else if user.Role.AccessLevel == 1 && userPermission.Permission.Resource == resource && userPermission.Permission.Scope == scope && userPermission.Permission.Group == group && userPermission.CreatedAt.Add(userPermission.Duration).Before(time.Now()) {
	// 		valid = true
	// 		break
	// 	} else if user.Role.AccessLevel > 1 {
	// 		valid = userPermission.Permission.Policy <= 3
	// 	}
	// }
}
