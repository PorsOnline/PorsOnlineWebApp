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

func (q *questionRepo) GetDB(ctx context.Context) *gorm.DB {
	return q.db
}

func (q *questionRepo) Create(ctx context.Context, question types.Question, tx *gorm.DB) (*types.Question, error) {
	err := tx.Model(&types.Question{}).Create(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (q *questionRepo) GetNextQuestionOrder(ctx context.Context, surveyID uint) (int, error) {
	var prevQuestion types.Question
	err := q.db.Order("\"order\" DESC").Where("is_dependency = false and survey_id = ?", surveyID).First(&prevQuestion).Error
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
	err := q.db.Model(&types.Question{}).Select("count(questions.id)>0").Joins("left join question_options op on questions.id = op.question_id").Where("op.next_question_id = ? and questions.deleted_at is null", id).Find(&dependencyExists).Error
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

func (q *questionRepo) Update(ctx context.Context, question types.Question, tx *gorm.DB) (*types.Question, error) {
	var oldQuestion types.Question
	err := q.db.Model(&types.Question{}).Where("id = ?", question.ID).First(&oldQuestion).Error
	if err != nil {
		return nil, err
	}
	err = tx.Model(&types.Question{}).Save(&question).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &question, nil
}

func (q *questionRepo) DeleteQuestionOptions(ctx context.Context, questionID uint, tx *gorm.DB) error { //todo: remove transaction from repo
	return tx.Model(&types.QuestionOption{}).Where("question_id = ?", questionID).Delete(&types.QuestionOption{}).Error
}

func (q *questionRepo) CreateQuestionOptions(ctx context.Context, options []types.QuestionOption, questionID uint, tx *gorm.DB) ([]types.QuestionOption, error) {
	for _, option := range options {
		err := tx.Model(&types.QuestionOption{}).Create(&types.QuestionOption{QuestionID: questionID, OptionText: option.OptionText, NextQuestionID: option.NextQuestionID}).Error
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

func (q *questionRepo) GetCurrentQuestion(ctx context.Context, userQuestionStep types.UserQuestionStep) (*types.UserQuestionStep, error) {
	var oldQuestion types.UserQuestionStep
	err := q.db.Model(&types.UserQuestionStep{}).
		Where("user_id = ? and survey_id = ?", userQuestionStep.UserID, userQuestionStep.SurveyID).
		Order("created_at desc").
		First(&oldQuestion).Error
	if err != nil {
		return nil, err
	}
	return &oldQuestion, nil
}

func (q *questionRepo) GetNextQuestionByCondition(ctx context.Context, userQuestionStep types.UserQuestionStep) (*uint, error) {
	var nextQuestionID uint
	err := q.db.Model(&types.Question{}).
		Select("op.next_question_id").
		Joins("left join question_options op on questions.id = op.question_id").
		Limit(1).
		Where("questions.survey_id = ? and questions.id = ? and op.next_question_id is not null", userQuestionStep.SurveyID, userQuestionStep.QuestionID).
		First(&nextQuestionID).Error
	
	// fmt.Println(*nextQuestionID, err)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &nextQuestionID, nil
}

func (q *questionRepo) GetNextQuestionByOrder(ctx context.Context, userQuestionStep types.UserQuestionStep) (*uint, error) {
	var questionID uint
	err := q.db.Model(&types.Question{}).
		Joins("left join user_question_steps uqs on uqs.question_id = questions.id ").
		Select("questions.id").
		Limit(1).
		Order("questions.\"order\" asc").
		Where("questions.survey_id = ? and uqs.id is null", userQuestionStep.SurveyID).
		First(&questionID).Error
	if err != nil {
		return nil, nil
	}
	return &questionID, nil
}

func (q *questionRepo) GetPreviousQuestion(ctx context.Context, userQuestionStep types.UserQuestionStep) (*uint, error) {
	var questionID *uint
	err := q.db.Model(&types.UserQuestionStep{}).
		Select("question_id").
		Limit(1).
		Order("created_at desc").
		Where("survey_id = ? ", userQuestionStep.SurveyID).
		Find(&questionID).Error
	if err != nil {
		return nil, err
	}
	return questionID, nil
}

func (q *questionRepo) GetFirstQuestion(ctx context.Context, surveyID uint) (*types.Question, error) {
	var question types.Question
	err := q.db.Model(&types.Question{}).Where("questions.survey_id = ? and \"order\" = 1 ", surveyID).First(&question).Error
	if err != nil {
		return &types.Question{}, err
	}
	return &question, nil
}

func (q *questionRepo) CreateQuestionStep(ctx context.Context, questionStep types.UserQuestionStep) error {
	return q.db.Model(&types.UserQuestionStep{}).Create(&questionStep).Error
}
