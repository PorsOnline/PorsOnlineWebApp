package domain

import (
	"github.com/porseOnline/pkg/adapters/storage/types"
)

type Question struct {
	ID              uint               `json:"id" validate:"omitempty"`
	SurveyID        uint               `json:"surveyUUID" validate:"required"`
	QuestionText    string             `json:"questionText" validate:"required"`
	IsDependency    bool               `json:"isDependency" validate:"omitempty"`
	QuestionType    types.QuestionType `json:"questionType" validate:"required,oneof=ConditionalMultipleChoice MultipleChoice MultipleChoiceWithAnswer Descriptive"`
	QuestionOptions []QuestionOption   `json:"questionOptions" validate:"omitempty"`
}

type QuestionOption struct {
	OptionText     string
	NextQuestionID *uint
	IsCorrect      bool
}

type UserQuestionStep struct {
	// SurveyUUID uuid.UUID    `json:"surveyUUID" validate:"required"`
	QuestionID uint         `json:"questionID" validate:"required"`
	Action     types.Action `json:"action" validate:"omitempty,oneof=forward backward"`
}

func TypeToDomainMapper(question types.Question) *Question {
	var questions []QuestionOption
	var isCorrect bool
	for _, option := range question.Options {
		if option.OptionText == question.CorrectAnswer {
			isCorrect = true
		} else {
			isCorrect = false
		}
		questions = append(questions, QuestionOption{OptionText: option.OptionText, NextQuestionID: option.NextQuestionID, IsCorrect: isCorrect})
	}

	return &Question{
		ID:              question.ID,
		SurveyID:        question.SurveyID,
		QuestionText:    question.QuestionText,
		IsDependency:    question.IsDependency,
		QuestionType:    question.QuestionType,
		QuestionOptions: questions,
	}
}

func DomainToTypeMapper(question Question) types.Question {
	var questions []types.QuestionOption
	var correctAnswer string
	for _, option := range question.QuestionOptions {
		if option.IsCorrect {
			correctAnswer = option.OptionText
		}
		questions = append(questions, types.QuestionOption{OptionText: option.OptionText, NextQuestionID: option.NextQuestionID})
	}
	return types.Question{
		SurveyID:      question.SurveyID,
		QuestionText:  question.QuestionText,
		IsDependency:  question.IsDependency,
		QuestionType:  question.QuestionType,
		Options:       questions,
		CorrectAnswer: correctAnswer,
	}
}

func QuestionStepDomainToType(question UserQuestionStep, surveyID uint) *types.UserQuestionStep {
	return &types.UserQuestionStep{
		QuestionID: question.QuestionID,
		SurveyID:   surveyID,
		Action:     question.Action,
	}
}
