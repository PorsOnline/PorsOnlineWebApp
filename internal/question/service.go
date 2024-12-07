package question

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/porseOnline/internal/question/domain"
	questionPort "github.com/porseOnline/internal/question/port"
	surveyPort "github.com/porseOnline/internal/survey/port"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"gorm.io/gorm"
)

type questionService struct {
	questionRepo  questionPort.Repo
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
	options, err := qs.questionRepo.CreateQuestionOptions(ctx, questionType.Options, createdQuestion.ID, tx)
	if err != nil {
		tx.Rollback()
		return domain.Question{}, err
	}
	tx.Commit()
	createdQuestion.Options = options
	return *domain.TypeToDomainMapper(*createdQuestion, survey.UUID), nil
}

func (qs *questionService) DeleteQuestion(ctx context.Context, id uint) error {
	err := qs.questionRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (qs *questionService) GetQuestionByID(ctx context.Context, id uint) (*domain.Question, error) {
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
	options, err := qs.questionRepo.CreateQuestionOptions(ctx, questionType.Options, updatedQuestion.ID, tx)
	if err != nil {
		tx.Rollback()
		return domain.Question{}, err
	}
	tx.Commit()
	updatedQuestion.Options = options
	return *domain.TypeToDomainMapper(*updatedQuestion, survey.UUID), nil
}

func (qs *questionService) validateUserInputsExistence(ctx context.Context, question domain.Question, surveyID uint) error {
	if question.QuestionType == types.ConditionalMultipleChoice {
		for _, option := range question.QuestionOptions {
			nextQuestionID, err := qs.questionRepo.Get(ctx, *option.NextQuestionID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New(fmt.Sprintf("next question for %v not found", option.OptionText))
				}
				return err
			}
			if nextQuestionID.SurveyID != surveyID {
				return errors.New("survey id mismatch")
			}
		}
	}
	return nil
}

func (qs questionService) GetNextQuestion(ctx context.Context, userQuestionStep domain.UserQuestionStep) (*domain.Question, error) {
	survey, err := qs.surveyService.GetSurveyByUUID(ctx, userQuestionStep.SurveyUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.Question{}, errors.New("survey not found")
		}
		return &domain.Question{}, err
	}
	questionStep := domain.QuestionStepDomainToType(userQuestionStep, survey.ID)
	userID, err := strconv.Atoi(ctx.Value("UserID").(string))
	questionStep.UserID = uint(userID)
	currentStep, err := qs.questionRepo.GetCurrentQuestion(ctx, *questionStep)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		question, err := qs.questionRepo.GetFirstQuestion(ctx, survey.ID)
		if err != nil {
			return &domain.Question{}, err
		}
		err = qs.questionRepo.CreateQuestionStep(ctx, types.UserQuestionStep{SurveyID: survey.ID, QuestionID: question.ID, UserID: uint(userID)})
		if err != nil {
			return &domain.Question{}, err
		}
		return domain.TypeToDomainMapper(*question, survey.UUID), nil
	} else if err != nil {
		return &domain.Question{}, err
	}
	if currentStep.QuestionID != questionStep.QuestionID {
		return &domain.Question{}, errors.New("not current step")
	}
	if questionStep.Action == types.Forward {
		nextQuestionID, err := qs.questionRepo.GetNextQuestionByCondition(ctx, *currentStep)
		if err != nil {
			return &domain.Question{}, err
		}
		if nextQuestionID == nil {
			nextQuestionID, err = qs.questionRepo.GetNextQuestionByOrder(ctx, *currentStep)
			if err != nil {
				return &domain.Question{}, err
			}
		}

		question, err := qs.questionRepo.Get(ctx, *nextQuestionID)
		err = qs.questionRepo.CreateQuestionStep(ctx, types.UserQuestionStep{SurveyID: survey.ID, QuestionID: question.ID, UserID: uint(userID), Action: questionStep.Action})
		if err != nil {
			return &domain.Question{}, err
		}
		return domain.TypeToDomainMapper(*question, survey.UUID), nil

	}
	previousQuestionID, err := qs.questionRepo.GetPreviousQuestion(ctx, *currentStep)
	if err != nil {
		return &domain.Question{}, err
	}
	question, err := qs.questionRepo.Get(ctx, *previousQuestionID)
	err = qs.questionRepo.CreateQuestionStep(ctx, types.UserQuestionStep{SurveyID: survey.ID, QuestionID: question.ID, UserID: uint(userID), Action: questionStep.Action})
	if err != nil {
		return &domain.Question{}, err
	}
	return domain.TypeToDomainMapper(*question, survey.UUID), nil

}
