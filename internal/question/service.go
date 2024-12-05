package question

import (
	"context"
	"errors"

	"github.com/porseOnline/internal/question/domain"
	questionPort "github.com/porseOnline/internal/question/port"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"gorm.io/gorm"
)

type questionService struct {
	questionRepo questionPort.Repo
	surveyService surveyPort.Service
}

func NewService(repo questionPort.Repo, surveyService surveyPort.Service) questionPort.Service {
	return &questionService{questionRepo: repo, surveyService: surveyService}
}

func (qs *questionService) CreateQuestion(ctx context.Context, question domain.Question) (domain.Question, error) {
	survey, err := qs.surveyService.GetSurveyByUUID(ctx, question.SurveyUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Question{}, errors.New("survey not found")
		}
		return domain.Question{}, err
	}
	err = qs.validateUserInputsExistence(ctx, question, survey.ID)
	if err != nil {
		return domain.Question{}, err
	}
	questionType := domain.DomainToTypeMapper(question, survey.ID)
	if !question.IsDependency {
		order, err := qs.questionRepo.GetNextQuestionOrder(ctx, questionType.SurveyID)
		if err != nil {
			return domain.Question{}, err
		}
		questionType.Order = order
	}
	tx := qs.questionRepo.GetDB(ctx).Begin()
	createdQuestion, err := qs.questionRepo.Create(ctx, questionType, tx)
	if err != nil {
		tx.Rollback()
		return domain.Question{}, err
	}
	options, err := qs.questionRepo.CreateQuestionOptions(ctx, questionType.Options, question.ID, tx)
	if err != nil {
		tx.Rollback()
		return domain.Question{}, err
	}
	tx.Commit()
	createdQuestion.Options = options
	return *domain.TypeToDomainMapper(*createdQuestion, survey.UUID), nil
}

func (qs *questionService) DeleteQuestion(ctx context.Context, id uint) (error) {
	err := qs.questionRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (qs *questionService) GetQuestion(ctx context.Context, id uint) (*domain.Question, error) {
	question, err := qs.questionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	survey, err := qs.surveyService.GetSurveyByID(ctx, question.ID)
	if err != nil {
		return nil, err
	}
	return domain.TypeToDomainMapper(*question, survey.UUID), nil
}

func (qs *questionService) UpdateQuestion(ctx context.Context, question domain.Question) (domain.Question, error) {
	survey, err := qs.surveyService.GetSurveyByUUID(ctx, question.SurveyUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Question{}, errors.New("survey not found")
		}
		return domain.Question{}, err
	}
	err = qs.validateUserInputsExistence(ctx, question, survey.ID)
	if err != nil {
		return domain.Question{}, err
	}
	questionType := domain.DomainToTypeMapper(question, survey.ID)
	tx := qs.questionRepo.GetDB(ctx).Begin()
	updatedQuestion, err := qs.questionRepo.Update(ctx, questionType, tx)
	if err != nil {
		tx.Rollback()
		return domain.Question{}, err
	}
	err = qs.questionRepo.DeleteQuestionOptions(ctx, question.ID, tx)
	if err != nil {
		tx.Rollback()
		return domain.Question{}, err
	}
	options, err := qs.questionRepo.CreateQuestionOptions(ctx, questionType.Options, question.ID, tx)
	if err != nil {
		tx.Rollback()
		return domain.Question{}, err
	}
	tx.Commit()
	updatedQuestion.Options = options
	return *domain.TypeToDomainMapper(*updatedQuestion, survey.UUID), nil
}

func (qs *questionService) validateUserInputsExistence(ctx context.Context, question domain.Question, surveyID uint) error {
	if question.NextQuestionIfFalseID != nil && question.NextQuestionIfTrueID != nil {
		nextQuestionIfFalseID, err := qs.questionRepo.Get(ctx, *question.NextQuestionIfFalseID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("next question if false not found")
			}
			return err
		}
		if nextQuestionIfFalseID.SurveyID != surveyID {
			return errors.New("survey id mismatch")
		}
		nextQuestionIfTrueID, err := qs.questionRepo.Get(ctx, *question.NextQuestionIfTrueID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("next question if true not found")
			}
			return err
		}
		if nextQuestionIfTrueID.SurveyID != surveyID {
			return errors.New("survey id mismatch")
		}
	}
	return nil
}
