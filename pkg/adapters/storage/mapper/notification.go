package mapper

import (
	"PorsOnlineWebApp/internal/notification/domain"
	"PorsOnlineWebApp/pkg/adapters/storage/types"
	"time"
)

func NotifDomain2Storage(notifDomain domain.Notification) *types.Notification {
	return &types.Notification{
		ID:        notifDomain.ID,
		UserID:    notifDomain.UserID,
		Message:   notifDomain.Message,
		Read:      false,
		Create_at: time.Now(),
	}
}

func NotifStorage2Domain(notifStorage []types.Notification) []*domain.Notification {
	result := make([]*domain.Notification, len(notifStorage))

	for i, notif := range notifStorage {
		result[i] = &domain.Notification{
			ID:        notif.ID,
			UserID:    notif.UserID,
			Message:   notif.Message,
			Read:      notif.Read,    
			Create_at: notif.Create_at,
		}
	}

	return result
}
