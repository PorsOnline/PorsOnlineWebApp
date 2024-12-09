package domain

import "time"

type (
	RoleID          uint
	TypeAccessLevel uint8
)

const (
	AccessLevelUnknown TypeAccessLevel = iota
	AccessLevelClient
	AccessLevelAdmin
	AccessLevelSuperAdmin
)

type Role struct {
	ID          RoleID
	CreatedAt   time.Time
	DeletedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	AccessLevel TypeAccessLevel
}
