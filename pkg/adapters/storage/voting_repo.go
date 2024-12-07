package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/porseOnline/internal/voting/domain"
	votingPort "github.com/porseOnline/internal/voting/port"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"github.com/porseOnline/pkg/helper"
	"github.com/porseOnline/pkg/logger"
	"gorm.io/gorm"
)

type submitRepo struct {
	db       *gorm.DB
	secretDB *gorm.DB
}

func NewVotingRepo(db *gorm.DB, secretDB *gorm.DB) votingPort.Repo {
	return &submitRepo{
		db:       db,
		secretDB: secretDB,
	}

}

func (su *submitRepo) Vote(ctx context.Context, answer *domain.Vote) error {
	if answer == nil {
		return errors.New("vote answer cannot be nil")
	}

	var oldSecret types.Secrets
	var key string

	err := su.secretDB.Table("secrets").
		WithContext(ctx).
		Where("user_id = ? AND servey_id = ?", answer.UserID, answer.SurveyID).
		First(&oldSecret).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		key = uuid.Must(uuid.NewUUID()).String()
		secret := types.Secrets{
			UserID:   answer.UserID,
			ServeyID: answer.SurveyID,
			Secret:   key,
		}

		if err := su.secretDB.Table("secrets").WithContext(ctx).Create(&secret).Error; err != nil {
			return err
		}
	} else {
		key = oldSecret.Secret
	}

	if answer.TextResponse != "" {
		cipherText, err := helper.EncryptAES(answer.TextResponse, []byte(key))
		if err != nil {
			logger.Error("failed to encrypt text response", nil)
			return err
		}
		answer.TextResponse = cipherText
	}

	if answer.SelectedOption != "" {
		cipherText, err := helper.EncryptAES(answer.SelectedOption, []byte(key))
		if err != nil {
			logger.Error("failed to encrypt selected option", nil)
			return err
		}
		answer.SelectedOption = cipherText
	}

	if err := su.db.Table("votes").WithContext(ctx).Create(answer).Error; err != nil {
		return err
	}

	return nil
}
