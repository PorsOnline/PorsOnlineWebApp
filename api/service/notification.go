package service

import (
	"PorsOnlineWebApp/internal/notification/domain"
	notificatinPort "PorsOnlineWebApp/internal/notification/port"
	"PorsOnlineWebApp/pkg/adapters/storage/types"
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type NotificationService struct {
	srv                   notificatinPort.Service
	authSecret            string
	expMin, refreshExpMin uint
}

func NewNotificationSerivce(srv notificatinPort.Service, authSecret string, expMin, refreshExpMin uint) *NotificationService {
	return &NotificationService{
		srv:           srv,
		authSecret:    authSecret,
		expMin:        expMin,
		refreshExpMin: refreshExpMin,
	}
}

func (n *NotificationService) GetUnreadMessages(ctx context.Context, userID string) ([]*domain.Notification, error) {
	var unreadMessages []*domain.Notification
	unreadMessages, err := n.srv.GetUnreadMessages(ctx, userID)
	if err != nil {
		log.Error("can not get unread messages")
		return unreadMessages, nil
	}
	return unreadMessages, nil
}

func (n *NotificationService) SendMessage(ctx context.Context, notif *types.Notification) error {
	err := n.srv.SendMessage(ctx, domain.Notification{
		ID:        notif.ID,
		UserID:    notif.UserID,
		Message:   notif.Message,
		Read:      false,
		Create_at: time.Now(),
	})
	if err != nil {
		log.Error("can not send message")
		return err
	}
	return nil
}
