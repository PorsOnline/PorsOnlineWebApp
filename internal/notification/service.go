package notification

import (
	"context"

	"github.com/porseOnline/internal/notification/domain"
	"github.com/porseOnline/internal/notification/port"

	"github.com/gofiber/fiber/v2/log"
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) SendMessage(ctx context.Context, notif domain.Notification) error {
	err := s.repo.SendMessage(ctx, notif)
	if err != nil {
		log.Errorf("can not send notif to user with id %d", notif.UserID)
		return err
	}

	return nil
}

func (s *service) GetUnreadMessages(ctx context.Context, userID string) ([]*domain.Notification, error) {
	var UnreadMessages []*domain.Notification
	UnreadMessages, err := s.repo.GetUnreadMessages(ctx, userID)
	if err != nil {
		log.Error("can not read unread messages")
		return []*domain.Notification{}, err
	}
	return UnreadMessages, nil

}
