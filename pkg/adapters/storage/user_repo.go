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

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) port.Repo {
	return &userRepo{db}

}

func (r *userRepo) Create(ctx context.Context, userDomain domain.User) (domain.UserID, error) {
	user := mapper.UserDomain2Storage(userDomain)
	return domain.UserID(user.ID), r.db.Table("users").WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	var user types.User
	err := r.db.Debug().Table("users").
		Where("id = ?", userID).WithContext(ctx).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, nil
	}

	return mapper.UserStorage2Domain(user), nil
}
func (r *userRepo) GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	var user types.User
	err := r.db.Table("users").
		Where("email = ?", email).
		First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user.ID == 0 {
		return nil, nil
	}

	return mapper.UserStorage2Domain(user), nil
}

func (r *userRepo) GetByFilter(ctx context.Context, filter *domain.UserFilter) (*domain.User, error) {
	var user types.User

	q := r.db.Table("users").Debug().WithContext(ctx)

	if filter.ID > 0 {
		q = q.Where("id = ?", filter.ID)
	}

	if len(filter.Phone) > 0 {
		q = q.Where("phone = ?", filter.Phone)
	}

	err := q.First(&user).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user.ID == 0 {
		return nil, nil
	}

	return mapper.UserStorage2Domain(user), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user domain.User) error {
	var preUpdateUser types.User
	err := r.db.Model(&types.User{}).Where("id = ?", user.ID).First((&preUpdateUser)).Error
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}
	currentTime := time.Now()
	if currentTime.Sub(preUpdateUser.CreatedAt) > 24*time.Hour {
		return errors.New("can not update user due to limitation of update time")
	}
	updates := make(map[string]interface{})
	if user.FirstName != "" {
		updates["first_name"] = user.FirstName
	}

	if user.FirstName != "" {
		updates["last_name"] = user.LastName
	}

	if user.Phone != "" {
		updates["phone"] = user.Phone
	}

	if user.Email != "" {
		updates["email"] = user.Email
	}

	if user.NationalCode != "" {
		updates["national_code"] = user.NationalCode
	}

	if user.BirthDate != preUpdateUser.BirthDate {
		updates["birth_date"] = user.BirthDate
	}

	if user.City != "" {
		updates["city"] = user.City
	}

	if user.Gender != preUpdateUser.Gender {
		updates["gender"] = user.Gender
	}

	if user.SurveyLimitNumber != preUpdateUser.SurveyLimitNumber {
		updates["survey_limit_number"] = user.SurveyLimitNumber
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		logger.Error(tx.Error.Error(), nil)
		return tx.Error
	}

	// Update the user record
	if err := tx.Model(&types.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
		logger.Error(err.Error(), nil)
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *userRepo) DeleteByID(ctx context.Context, userID domain.UserID) error{
	return	r.db.Where("id = ?", userID).Delete(&types.User{}).Error
}