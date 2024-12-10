package types

import (
	codeVerifDomain "github.com/porseOnline/internal/codeVerification/domain"
	"github.com/porseOnline/internal/common"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OutboxData struct {
	Dest    string
	Content string
	Type    codeVerifDomain.CodeVerificationType
}

type CodeVerificationOutbox struct {
	gorm.Model
	OutboxID           common.OutboxID
	CodeVerificationID codeVerifDomain.CodeVerificationID
	Data               datatypes.JSON
	Status             common.OutboxStatus
	Type               common.OutboxType
}
type Outbox struct {
	gorm.Model
	Data   datatypes.JSON
	RefID  uint
	Type   uint8
	Status uint8
}

type CodeVerification struct {
	gorm.Model
	Content string
	To      uint
	Type    uint8
}
