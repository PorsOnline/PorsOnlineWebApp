package port

import (
	"context"

	"github.com/porseOnline/internal/codeVerification/domain"
	"github.com/porseOnline/internal/common"
	userDomain "github.com/porseOnline/internal/user/domain"
)

type Repo interface {
	Create(ctx context.Context, codeVerification *domain.CodeVerification) (domain.CodeVerificationID, error)
	CreateOutbox(ctx context.Context, outbox *domain.CodeVerificationOutbox) error
	QueryOutboxes(ctx context.Context, limit uint, status common.OutboxStatus) ([]domain.CodeVerificationOutbox, error)
	GetUserCodeVerificationValue(ctx context.Context, userID userDomain.UserID) (string, error)
}
