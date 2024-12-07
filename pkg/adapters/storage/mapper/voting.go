package mapper

import (
	"github.com/porseOnline/internal/voting/domain"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"gorm.io/gorm"
)

func VotingDomain2Storage(votingDomain domain.Vote) *types.Vote {
	return &types.Vote{
		Model: gorm.Model{
			ID:        votingDomain.ID,
			CreatedAt: votingDomain.CreatedAt,
			DeletedAt: gorm.DeletedAt(ToNullTime(votingDomain.DeletedAt)),
			UpdatedAt: votingDomain.UpdatedAt,
		},
		UserID:         votingDomain.UserID,
		SurveyID:       votingDomain.SurveyID,
		QuestionID:     votingDomain.QuestionID,
		TextResponse:   votingDomain.TextResponse,
		SelectedOption: votingDomain.SelectedOption,
	}
}

func VotingStorage2Domain(votingStorage []types.Vote) []*domain.Vote {
	result := make([]*domain.Vote, len(votingStorage))

	for i, voting := range votingStorage {
		result[i] = &domain.Vote{
			UserID:         voting.UserID,
			SurveyID:       voting.SurveyID,
			QuestionID:     voting.QuestionID,
			TextResponse:   voting.TextResponse,
			SelectedOption: voting.SelectedOption,
		}
	}

	return result
}
