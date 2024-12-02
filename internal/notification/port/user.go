package port

import (
	"PorsOnlineWebApp/internal/notification/domain"
	"context"
)

type Repo interface{
	SendMessage(ctx context.Context, notif domain.Notification) (error)
	GetUnreadMessages(ctx context.Context, userID string) ([]*domain.Notification, error)
}