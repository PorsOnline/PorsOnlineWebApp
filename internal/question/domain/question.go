package domain

import (
	"github.com/google/uuid"
	"github.com/porseOnline/pkg/adapters/storage/types"
)

type Question struct {
	ID                    uint               `json:"id" validate:"omitempty"`
	SurveyUUID            uuid.UUID          `json:"surveyUUID,required" validate:"required,uuid"`
	QuestionText          string             `json:"questionText" validate:"required"`
	NextQuestionIfTrueID  *uint              `json:"nextQuestionIfTrueID" validate:"omitempty"`
	NextQuestionIfFalseID *uint              `json:"nextQuestionIfFalseID" validate:"omitempty"`
	CorrectAnswer         string             `json:"correctAnswer" validate:"omitempty"`
	IsDependency          bool               `json:"isDependency" validate:"omitempty"`
	QuestionType          types.QuestionType `json:"questionType" validate:"omitempty,oneof=Conditional ConditionalWithAnswer MultipleChoice MultipleChoiceWithAnswer Descriptive"`
	QuestionOptions       []QuestionOption   `json:"questionOptions" validate:"required"`
}

type QuestionOption struct {
	OptionText string
}

func TypeToDomainMapper(question types.Question, surveyUUID uuid.UUID) *Question {
	var questions []QuestionOption
	for _, option := range question.Options {
		questions = append(questions, QuestionOption{OptionText: option.OptionText})
	}
	return &Question{
		ID:                    question.ID,
		SurveyUUID:            surveyUUID,
		QuestionText:          question.QuestionText,
		NextQuestionIfTrueID:  question.NextQuestionIfTrueID,
		NextQuestionIfFalseID: question.NextQuestionIfFalseID,
		IsDependency:          question.IsDependency,
		QuestionType:          question.QuestionType,
		QuestionOptions:       questions,
		CorrectAnswer:         question.CorrectAnswer,
	}
}

func DomainToTypeMapper(question Question, surveyID uint) types.Question {
	var questions []types.QuestionOption
	for _, option := range question.QuestionOptions {
		questions = append(questions, types.QuestionOption{OptionText: option.OptionText})
	}
	return types.Question{
		SurveyID:              surveyID,
		QuestionText:          question.QuestionText,
		NextQuestionIfTrueID:  question.NextQuestionIfTrueID,
		NextQuestionIfFalseID: question.NextQuestionIfFalseID,
		IsDependency:          question.IsDependency,
		QuestionType:          question.QuestionType,
		Options:               questions,
		CorrectAnswer:         question.CorrectAnswer,
	}
}
