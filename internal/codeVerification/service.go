package codeVerification

import (
	"context"
	"fmt"
	"time"

	"github.com/porseOnline/internal/codeVerification/domain"
	codeVerificationRepo "github.com/porseOnline/internal/codeVerification/port"
	"github.com/porseOnline/internal/common"
	userDomain "github.com/porseOnline/internal/user/domain"
	userPort "github.com/porseOnline/internal/user/port"

	"github.com/porseOnline/pkg/fp"
	"github.com/porseOnline/pkg/helper"
)

type Service struct {
	codeVerificationRepo codeVerificationRepo.Repo
	outboxRepo           common.OutboxRepo
	userPort             userPort.Service
}

func NewService(userPort userPort.Service, outboxRepo common.OutboxRepo, codeVerificationRepo codeVerificationRepo.Repo) codeVerificationRepo.Service {
	return &Service{
		userPort:             userPort,
		codeVerificationRepo: codeVerificationRepo,
		outboxRepo:           outboxRepo,
	}
}

func (s *Service) Send(ctx context.Context, codeVerification *domain.CodeVerification) error {

	user, err := s.userPort.GetUserByFilter(ctx, &userDomain.UserFilter{
		ID: codeVerification.UserID,
	})

	if err != nil {
		return err
	}

	codeVerificationID, err := s.codeVerificationRepo.Create(ctx, codeVerification)
	if err != nil {
		return err
	}
	err = s.codeVerificationRepo.CreateOutbox(ctx, &domain.CodeVerificationOutbox{
		CodeVerificationID: codeVerificationID,
		Data: domain.OutboxData{
			Dest: func() string {
				switch codeVerification.Type {
				case domain.CodeVerificationTypeSMS:
					return string(user.Phone)
				case domain.CodeVerificationTypeEmail:
					return string(user.Email)
				default:
					return ""
				}
			}(),
			Content: codeVerification.Content,
			Type:    codeVerification.Type,
		},
		Status: common.OutboxStatusCreated,
		Type:   common.OutboxTypeCodeVerification,
	})
	if err != nil {
		return err
	}
	/////////

	return nil
}

func (s *Service) Handle(ctx context.Context, outboxes []domain.CodeVerificationOutbox) error {
	outBoxIDs := fp.Map(outboxes, func(o domain.CodeVerificationOutbox) common.OutboxID {
		return o.OutboxID
	})

	if err := s.outboxRepo.UpdateBulkStatuses(ctx, common.OutboxStatusPicked, outBoxIDs...); err != nil {
		return fmt.Errorf("failed to update code outbox statuses to picked %w", err)
	}

	for _, outbox := range outboxes {
		// fmt.Printf("dest : %s, content : %s\n", outbox.Data.Dest, outbox.Data.Content)
		go helper.SendEmail(outbox.Data.Dest, outbox.Data.Content)

	}

	if err := s.outboxRepo.UpdateBulkStatuses(ctx, common.OutboxStatusDone, outBoxIDs...); err != nil {
		return fmt.Errorf("failed to update code outbox statuses to done %w", err)
	}

	return nil
}

func (s *Service) Interval() time.Duration {
	return time.Second * 10
}

func (s *Service) Query(ctx context.Context) ([]domain.CodeVerificationOutbox, error) {
	return s.codeVerificationRepo.QueryOutboxes(ctx, 100, common.OutboxStatusCreated)
}

func (s *Service) CheckUserCodeVerificationValue(ctx context.Context, userID userDomain.UserID, val string) (bool, error) {
	expected, err := s.codeVerificationRepo.GetUserCodeVerificationValue(ctx, userID)
	if err != nil {
		return false, err
	}
	return expected == val, nil
}
