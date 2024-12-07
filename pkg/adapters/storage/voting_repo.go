package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/porseOnline/internal/voting/domain"
	votingPort "github.com/porseOnline/internal/voting/port"
	"github.com/porseOnline/pkg/adapters/storage/mapper"
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
	storageAnswer := mapper.VotingDomain2Storage(*answer)

	var oldSecret types.Secrets
	var randomBytes []byte
	err := su.secretDB.Table("secrets").
		WithContext(ctx).
		Where("user_id = ? AND servey_id = ?", storageAnswer.UserID, storageAnswer.SurveyID).
		First(&oldSecret).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		randomBytes = make([]byte, 16)

		// Read random bytes into the slice
		_, err := rand.Read(randomBytes)
		if err != nil {
			return err
		}
		// key := hex.EncodeToString(randomBytes)

		secret := types.Secrets{
			UserID:   storageAnswer.UserID,
			ServeyID: storageAnswer.SurveyID,
			Secret:   hex.EncodeToString(randomBytes),
		}

		if err := su.secretDB.Table("secrets").WithContext(ctx).Create(&secret).Error; err != nil {
			return err
		}
	} else {
		randomBytes, _ = hex.DecodeString(oldSecret.Secret)
	}

	if storageAnswer.TextResponse != "" {
		cipherText, err := helper.EncryptAES(storageAnswer.TextResponse, randomBytes)
		if err != nil {
			logger.Error("failed to encrypt text response", nil)

			return err
		}
		storageAnswer.TextResponse = cipherText
	}

	if storageAnswer.SelectedOption != "" {
		cipherText, err := helper.EncryptAES(storageAnswer.SelectedOption, randomBytes)
		if err != nil {
			logger.Error("failed to encrypt selected option", nil)
			return err
		}
		storageAnswer.SelectedOption = cipherText
	}
	if err := su.db.Table("votes").WithContext(ctx).Create(storageAnswer).Error; err != nil {
		return err
	}

	return nil
}

func (su *submitRepo) GetLastAnswer(ctx context.Context, userID uint, surveyID uint) (domain.Vote, error) {
	var lastAnswer domain.Vote
	var secret types.Secrets
	err := su.secretDB.Table("secrets").Where("user_id = ? AND survey_id = ?", userID, surveyID).Find(&secret).Error
	if err != nil {
		return domain.Vote{}, err
	}
	err = su.db.Table("votes").Where("user_id = ? AND survey_id = ?", userID, surveyID).Order("question_id desc").Find(&lastAnswer).Error
	if err != nil {
		return domain.Vote{}, err
	}
	key, err := hex.DecodeString(secret.Secret)
	if err != nil {
		return domain.Vote{}, err
	}
	if lastAnswer.TextResponse != "" {

		lastAnswer.TextResponse, err = helper.DecryptAES(lastAnswer.TextResponse, key)
		if err != nil {
			return domain.Vote{}, err
		}
	}

	if lastAnswer.SelectedOption != "" {
		lastAnswer.SelectedOption, err = helper.DecryptAES(lastAnswer.TextResponse, key)
		if err != nil {
			return domain.Vote{}, err
		}
	}
	return lastAnswer, nil

}
