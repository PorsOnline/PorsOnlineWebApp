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
	err := q.db.Debug().Order("\"order\" DESC").Where("is_dependency = false and survey_id = ?", surveyID).First(&prevQuestion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1, nil
		} else {
			return 0, err
		}
	}
	return prevQuestion.Order + 1, nil
}

func (q *questionRepo) Delete(ctx context.Context, id uint) error {
	var dependencyExists bool
	err := q.db.Model(&types.Question{}).Select("count(id)>0").Where("(next_question_if_true_id = ? or next_question_if_false_id = ?) and deleted_at is null", id, id).Find(&dependencyExists).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if dependencyExists {
		return errors.New("can not delete question")
	}
	var question types.Question
	err = q.db.Model(&types.Question{}).Where("id = ?", id).First(&question).Error
	if err != nil {
		return err
	}
	err = q.db.Model(&types.Question{}).Where("id = ?", id).Delete(&question).Error
	if err != nil {
		return err
	}
	return nil
}

func (q *questionRepo) Get(ctx context.Context, id uint) (*types.Question, error) {
	var question types.Question
	err := q.db.Model(&types.Question{}).Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (q *questionRepo) Update(ctx context.Context, question types.Question) (*types.Question, error) {
	var oldQuestion types.Question
	err := q.db.Debug().Model(&types.Question{}).Where("id = ?", question.ID).First(&oldQuestion).Error
	if err != nil {
		return nil, err
	}
	tx := q.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	err = tx.Model(&types.QuestionOption{}).Where("question_id = ?", question.ID).Delete(&types.QuestionOption{}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Debug().Model(&types.Question{}).Save(&question).Error
	if err != nil {
		tx.Rollback()
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
	
