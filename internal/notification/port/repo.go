package port

import (
	"context"

	"github.com/porseOnline/internal/notification/domain"
)

type Repo interface {
	SendMessage(ctx context.Context, notif domain.Notification) error
	GetUnreadMessages(ctx context.Context, userID string) ([]*domain.Notification, error)
}
