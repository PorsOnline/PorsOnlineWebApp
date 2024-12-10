package mapper

import (
	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

func UserDomain2Storage(userDomain domain.User) *types.User {
	var userRole types.Role
	if userRole.ID > 0 {
		userRole = *RoleDomain2Storage(userDomain.Role)
	} else {
		userRole = types.Role{}
	}

	return &types.User{
		Model: gorm.Model{
			ID:        uint(userDomain.ID),
			CreatedAt: userDomain.CreatedAt,
			DeletedAt: gorm.DeletedAt(ToNullTime(userDomain.DeletedAt)),
			UpdatedAt: userDomain.UpdatedAt,
		},
		FirstName:         userDomain.FirstName,
		LastName:          userDomain.LastName,
		Phone:             string(userDomain.Phone),
		Email:             string(userDomain.Email),
		PasswordHash:      userDomain.PasswordHash,
		NationalCode:      string(userDomain.NationalCode),
		BirthDate:         userDomain.BirthDate,
		City:              userDomain.City,
		Gender:            userDomain.Gender,
		SurveyLimitNumber: userDomain.SurveyLimitNumber,
		Balance:           userDomain.Balance,
		Role:              &userRole,
		RoleID:            (*uint)(&userDomain.Role.ID),
	}
}

func UserStorage2Domain(user types.User) *domain.User {
	var userPermissions []domain.Permission
	if len(user.UserPermissions) > 0 {
		for _, userPermission := range user.UserPermissions {
			userPermissions = append(userPermissions, *PermissionStorage2Domain(userPermission.Permission))
		}
	}

	var userRole domain.Role
	if userRole.ID > 0 {
		userRole = *RoleStorage2Domain(*user.Role)
	} else {
		userRole = domain.Role{}
	}

	return &domain.User{
		ID:        domain.UserID(user.ID),
		CreatedAt: user.CreatedAt,
		// DeletedAt:         user.DeletedAt,
		UpdatedAt:         user.UpdatedAt,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Phone:             domain.Phone(user.Phone),
		Email:             domain.Email(user.Email),
		PasswordHash:      user.PasswordHash,
		NationalCode:      domain.NationalCode(user.NationalCode),
		BirthDate:         user.BirthDate,
		City:              user.City,
		Gender:            user.Gender,
		SurveyLimitNumber: user.SurveyLimitNumber,
		Balance:           user.Balance,
		Role:              userRole,
		Permissions:       userPermissions,
	}
}
