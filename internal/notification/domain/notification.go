package domain

import "time"

type Notification struct {
	ID        string
	UserID    string
	Message   string
	Read      bool
	Create_at time.Time
}
