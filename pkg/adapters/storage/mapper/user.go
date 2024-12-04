package mapper

import (
	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

func UserDomain2Storage(userDomain domain.User) *types.User {
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
	}
}

func UserStorage2Domain(user types.User) *domain.User {
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
	}
}
