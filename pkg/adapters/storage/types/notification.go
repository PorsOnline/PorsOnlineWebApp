package types

import "time"

type Notification struct {
	ID        string    `json:"id,omitempty"  yaml:"id"`
	UserID    string    `json:"user_id"  yaml:"user_id"`
	Message   string    `json:"message"  yaml:"message"`
	Read      bool      `json:"read,omitempty"  yaml:"read"`
	Create_at time.Time `json:"create_at,omitempty"  yaml:"create_at"`
}
