package mapper

import (
	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

func PermissionDomain2Storage(permissionDomain domain.Permission) *types.Permission {
	return &types.Permission{
		Model: gorm.Model{
			ID:        uint(permissionDomain.ID),
			CreatedAt: permissionDomain.CreatedAt,
			DeletedAt: gorm.DeletedAt(ToNullTime(permissionDomain.DeletedAt)),
			UpdatedAt: permissionDomain.UpdatedAt,
		},
		Owner:    uint(permissionDomain.Owner),
		Group:    permissionDomain.Group,
		Resource: permissionDomain.Resource,
		Scope:    permissionDomain.Scope,
		Policy:   uint8(permissionDomain.Policy),
	}
}

func PermissionStorage2Domain(permission types.Permission) *domain.Permission {
	var users []domain.User
	for _, userPermission := range permission.UserPermissions {
		users = append(users, *UserStorage2Domain(*userPermission.User))
	}

	return &domain.Permission{
		ID:        domain.PermissionID(permission.ID),
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
		Owner:     domain.OwnerID(permission.Owner),
		Group:     permission.Group,
		Resource:  permission.Resource,
		Scope:     permission.Scope,
		Policy:    domain.TypePolicy(permission.Policy),
		// Duration:  permission.Duration,
		Users: users,
	}
}

func PermissionDetailsDomain2Storage(permissionDetailsDomain domain.PermissionDetails) *types.UserPermission {
	return &types.UserPermission{
		PermissionID: permissionDetailsDomain.PermissionID,
		SurveyID:     permissionDetailsDomain.SurveyID,
		Duration:     *permissionDetailsDomain.Duration,
	}
}
