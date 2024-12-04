package storage

import (
	"context"
	"errors"

	"github.com/porseOnline/internal/question/port"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"gorm.io/gorm"
)

type questionRepo struct {
	db *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) port.Repo {
	return &questionRepo{db: db}
}

func (q *questionRepo) Create(ctx context.Context, question types.Question) (*types.Question, error) {
	tx := q.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	err := tx.Debug().Model(&types.Question{}).Create(&question).Error
	if err != nil {
		return nil, err
	}
	for _, option := range question.Options {
		err := tx.Model(&types.QuestionOption{}).
			Create(&types.QuestionOption{
				QuestionID: question.ID,
				OptionText: option.OptionText,
				IsCorrect:  option.IsCorrect,
			}).
			Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &question, nil
}

func (q *questionRepo) GetNextQuestionOrder(ctx context.Context, surveyID uint) (int, error) {
	var prevQuestion types.Question
	err := q.db.Model(&types.Question{}).Where("is_dependency is null, survey_id = ?", surveyID).Order("order desc").First(&prevQuestion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1, nil
		} else {
			return 0, err
		}
	}
	return prevQuestion.Order + 1, nil
}