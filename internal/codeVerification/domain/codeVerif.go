package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/porseOnline/internal/common"
	userDomain "github.com/porseOnline/internal/user/domain"
)

type (
	CodeVerificationID     uint
	CodeVerificationType   uint8
	CodeVerificationStatus uint8
)

const (
	CodeVerificationTypeSMS CodeVerificationType = iota + 1
	CodeVerificationTypeEmail
)

const (
	CodeVerificationStatusCreated CodeVerificationStatus = iota + 1
	CodeVerificationStatusSent
)

type CodeVerification struct {
	ID            CodeVerificationID
	CreatedAt     time.Time
	UserID        userDomain.UserID
	Type          CodeVerificationType
	Content       string
	ForValidation bool
	TTL           time.Duration
}

func (n *CodeVerification) Normalize() {
	n.Content = strings.TrimSpace(n.Content)
	n.ID = 0
}

func (n *CodeVerification) Validate() error {
	if n.UserID == 0 {
		return errors.New("empty user id")
	}

	return nil
}

func NewCodeVerification(userID userDomain.UserID, content string, CodeVerificationType CodeVerificationType, forValidation bool, ttl time.Duration) *CodeVerification {
	return &CodeVerification{
		UserID:        userID,
		Type:          CodeVerificationType,
		Content:       content,
		CreatedAt:     time.Now(),
		ForValidation: forValidation,
		TTL:           ttl,
	}
}

type OutboxData struct {
	Dest    string
	Content string
	Type    CodeVerificationType
}

type CodeVerificationOutbox struct {
	OutboxID           common.OutboxID
	CodeVerificationID CodeVerificationID
	Data               OutboxData
	Status             common.OutboxStatus
	Type               common.OutboxType
}
