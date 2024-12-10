package port

import (
	"context"

	"github.com/porseOnline/internal/codeVerification/domain"
	"github.com/porseOnline/internal/common"
	userDomain "github.com/porseOnline/internal/user/domain"
)

type Service interface {
	Send(ctx context.Context, codeVerification *domain.CodeVerification) error
	CheckUserCodeVerificationValue(ctx context.Context, userID userDomain.UserID, val string) (bool, error)
	common.OutboxHandler[domain.CodeVerificationOutbox]
}
