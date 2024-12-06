package storage

import (
	"context"
	"errors"

	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/internal/user/port"
	"github.com/porseOnline/pkg/adapters/storage/mapper"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"github.com/porseOnline/pkg/logger"
	"gorm.io/gorm"
)

type roleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) port.RoleRepo {
	return &roleRepo{db}

}

func (r *roleRepo) Create(ctx context.Context, roleDomain domain.Role) (domain.RoleID, error) {
	role := mapper.RoleDomain2Storage(roleDomain)
	return domain.RoleID(role.ID), r.db.Table("roles").WithContext(ctx).Create(role).Error
}

func (r *roleRepo) Get(ctx context.Context, roleID domain.RoleID) (*domain.Role, error) {
	var role types.Role
	err := r.db.Debug().Table("roles").Where("id = ?", roleID).WithContext(ctx).First(&role).Error

	if err != nil {
		return nil, err
	}

	if role.ID == 0 {
		return nil, errors.New("role not found")
	}

	return mapper.RoleStorage2Domain(role), nil
}

func (r *roleRepo) Update(ctx context.Context, roleDomain domain.Role) error {
	var updatingRole types.Role
	err := r.db.Model(&types.Role{}).Where("id = ?", roleDomain.ID).First((&updatingRole)).Error
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}

	updates := make(map[string]interface{})
	if roleDomain.Name != "" {
		updates["name"] = roleDomain.Name
	}

	if roleDomain.AccessLevel != 0 {
		updates["access_level"] = roleDomain.AccessLevel
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Error(tx.Error.Error(), nil)
		return tx.Error
	}

	if err := tx.Model(&types.Role{}).Where("id = ?", roleDomain.ID).Updates(updates).Error; err != nil {
		logger.Error(err.Error(), nil)
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *roleRepo) Delete(ctx context.Context, roleID domain.RoleID) error {
	return r.db.Where("id = ?", roleID).Delete(&types.Role{}).Error
}

func (r *roleRepo) Assign(ctx context.Context, roleID domain.RoleID, userID domain.UserID) error {
	var role types.Role
	err := r.db.Debug().Table("roles").Where("id = ?", roleID).WithContext(ctx).First(&role).Error

	if err != nil {
		return err
	}

	if role.ID == 0 {
		return errors.New("role not found")
	}

	// find user function
	var user types.User
	err = r.db.Debug().Table("users").Where("id = ?", userID).WithContext(ctx).First(&user).Error

	if err != nil {
		return err
	}

	if user.ID == 0 {
		return errors.New("user not found")
	}

	//update user role
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Error(tx.Error.Error(), nil)
		return tx.Error
	}

	user.RoleID = role.ID
	user.Role = role
	if err := tx.Model(&types.User{}).Where("id = ?", userID).Updates(user).Error; err != nil {
		logger.Error(err.Error(), nil)
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
