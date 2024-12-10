package domain

import "time"

type (
	TypePolicy   uint8
	PermissionID uint
	OwnerID      uint
)

const (
	PolicyUnknown TypePolicy = iota
	PolicyClient
	PolicyAdmin
	PolicySuperAdmin
	PolicyOwner
)

type Permission struct {
	ID        PermissionID
	CreatedAt time.Time
	DeletedAt time.Time
	UpdatedAt time.Time
	Owner     OwnerID
	Group     string
	Resource  string
	Scope     string
	Policy    TypePolicy
	Users     []User
}
