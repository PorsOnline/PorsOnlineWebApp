package mapper

import (
	"encoding/json"

	"github.com/porseOnline/internal/codeVerification/domain"
	"github.com/porseOnline/internal/common"
	"github.com/porseOnline/pkg/adapters/storage/types"
)

func CodeVerificationOutbox2Storage(no *domain.CodeVerificationOutbox) (*types.Outbox, error) {
	data, err := json.Marshal(&no.Data)
	if err != nil {
		return nil, err
	}

	return &types.Outbox{
		Data:   data,
		RefID:  uint(no.CodeVerificationID),
		Type:   uint8(no.Type),
		Status: uint8(no.Status),
	}, nil
}

func CodeVerification2Storage(cv *domain.CodeVerification) *types.CodeVerification {
	return &types.CodeVerification{
		Content: cv.Content,
		To:      uint(cv.UserID),
		Type:    uint8(cv.Type),
	}
}

func OutboxStorage2CodeVerification(outbox types.Outbox) (domain.CodeVerificationOutbox, error) {
	var outboxData domain.OutboxData
	err := json.Unmarshal([]byte(outbox.Data), &outboxData)
	if err != nil {
		return domain.CodeVerificationOutbox{}, err
	}

	return domain.CodeVerificationOutbox{
		OutboxID:           common.OutboxID(outbox.ID),
		CodeVerificationID: domain.CodeVerificationID(outbox.RefID),
		Data:               outboxData,
		Status:             common.OutboxStatus(outbox.Status),
		Type:               common.OutboxType(outbox.Type),
	}, nil
}
