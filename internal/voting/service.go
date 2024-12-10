package voting

import (
	"context"

	"github.com/porseOnline/internal/voting/domain"
	"github.com/porseOnline/internal/voting/port"
	"github.com/porseOnline/pkg/logger"
)

type service struct {
	repo port.Repo
}

func NewVotingService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Vote(ctx context.Context, answer *domain.Vote) error {
	err := s.repo.Vote(ctx, answer)
	if err != nil {
		logger.Error("Can not vote", nil)
		return err
	}
	return nil
}

func (s *service) GetLastResponse(ctx context.Context, userID uint, serveyID uint) (domain.Vote, error) {
	lastAnswer, err := s.repo.GetLastResponse(ctx, userID, serveyID)
	if err!=nil{
		return domain.Vote{}, err
	}
	return lastAnswer, nil
}
