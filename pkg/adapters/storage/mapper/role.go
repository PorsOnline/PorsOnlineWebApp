package mapper

import (
	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

func RoleDomain2Storage(roleDomain domain.Role) *types.Role {
	return &types.Role{
		Model: gorm.Model{
			ID:        uint(roleDomain.ID),
			CreatedAt: roleDomain.CreatedAt,
			DeletedAt: gorm.DeletedAt(ToNullTime(roleDomain.DeletedAt)),
			UpdatedAt: roleDomain.UpdatedAt,
		},
		Name:        roleDomain.Name,
		AccessLevel: uint8(roleDomain.AccessLevel),
	}
}

func RoleStorage2Domain(role types.Role) *domain.Role {
	return &domain.Role{
		ID:          domain.RoleID(role.ID),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
		Name:        role.Name,
		AccessLevel: domain.TypeAccessLevel(role.AccessLevel),
	}
}
