package storage

import (
	"context"
	"errors"

	"github.com/porseOnline/internal/codeVerification/domain"
	"github.com/porseOnline/internal/codeVerification/port"
	"github.com/porseOnline/internal/common"
	userDomain "github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/pkg/adapters/storage/mapper"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"gorm.io/gorm"
)

type CodeVerificationRepo struct {
	db *gorm.DB
}

func NewCodeVerificationRepo(db *gorm.DB) port.Repo {
	return &CodeVerificationRepo{
		db: db,
	}
}

func (r *CodeVerificationRepo) Create(ctx context.Context, codeverification *domain.CodeVerification) (domain.CodeVerificationID, error) {
	no := mapper.CodeVerification2Storage(codeverification)

	if err := r.db.WithContext(ctx).Table("code_verifications").Create(no).Error; err != nil {
		return 0, err
	}

	return domain.CodeVerificationID(no.ID), nil
}

func (r *CodeVerificationRepo) CreateOutbox(ctx context.Context, no *domain.CodeVerificationOutbox) error {
	outbox, err := mapper.CodeVerificationOutbox2Storage(no)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Table("outboxes").Create(outbox).Error
}

func (r *CodeVerificationRepo) QueryOutboxes(ctx context.Context, limit uint, status common.OutboxStatus) ([]domain.CodeVerificationOutbox, error) {
	var outboxes []types.Outbox

	err := r.db.WithContext(ctx).Table("outboxes").
		Where(`"type" = ?`, common.OutboxTypeCodeVerification).
		Where("status = ?", status).
		Limit(int(limit)).Scan(&outboxes).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	result := make([]domain.CodeVerificationOutbox, len(outboxes))

	for i := range outboxes {
		v, err := mapper.OutboxStorage2CodeVerification(outboxes[i])
		if err != nil {
			return nil, err
		}
		result[i] = v
	}

	return result, nil
}

func (r *CodeVerificationRepo) GetUserCodeVerificationValue(ctx context.Context, userID userDomain.UserID) (string, error) {
	var code string
	err := r.db.WithContext(ctx).
		Table("code_verifications").
		Select("content").         // Get "content" from the table
		Where(`"to" = ?`, userID). // Use simple column name
		Scan(&code).Error

	// Check if the error is "record not found"
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil // Return an empty value instead of an error
	}

	return code, err

}
