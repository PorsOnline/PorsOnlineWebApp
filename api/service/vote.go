package service

import (
	"context"

	"github.com/porseOnline/internal/voting/domain"
	votingPort "github.com/porseOnline/internal/voting/port"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"github.com/porseOnline/pkg/logger"
)

type VoteService struct {
	srv                   votingPort.Service
	authSecret            string
	expMin, refreshExpMin uint
}

func NewVotingService(srv votingPort.Service, authSecret string, expMin, refreshExpMin uint) *VoteService {
	return &VoteService{
		srv:           srv,
		authSecret:    authSecret,
		expMin:        expMin,
		refreshExpMin: refreshExpMin,
	}
}

func (v *VoteService) Vote(ctx context.Context, answer *types.Vote) error {
	err := v.srv.Vote(ctx, &domain.Vote{
		UserID:         answer.UserID,
		SurveyID:       answer.SurveyID,
		QuestionID:     answer.QuestionID,
		TextResponse:   answer.TextResponse,
		SelectedOption: answer.SelectedOption,
	})
	if err != nil {
		logger.Error("can not vote", nil)
		return err
	}
	return err
}
