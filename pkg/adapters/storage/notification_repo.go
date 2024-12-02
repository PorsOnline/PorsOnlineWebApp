package storage

import (
	"context"
	"errors"

	"github.com/porseOnline/internal/notification/domain"
	notifPort "github.com/porseOnline/internal/notification/port"
	"github.com/porseOnline/pkg/adapters/storage/mapper"
	"github.com/porseOnline/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

type notifRepo struct {
	db *gorm.DB
}

func NewNotifRepo(db *gorm.DB) notifPort.Repo {
	return &notifRepo{
		db: db,
	}
}

func (r *notifRepo) SendMessage(ctx context.Context, notif domain.Notification) error {
	o := mapper.NotifDomain2Storage(notif)
	return r.db.Table("inbox").WithContext(ctx).Create(o).Error
}

func (r *notifRepo) GetUnreadMessages(ctx context.Context, userID string) ([]*domain.Notification, error) {
	var notif []types.Notification

	err := r.db.Table("inbox").WithContext(ctx).Where("user_id = ? AND read = ?", userID, false).Find(&notif).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return mapper.NotifStorage2Domain(notif), nil
}
