package domain

import (
	"github.com/google/uuid"
	"github.com/porseOnline/pkg/adapters/storage/types"
)

type Question struct {
	ID                    uint               `json:"id" validate:"omitempty"`
	SurveyUUID            uuid.UUID          `json:"surveyUUID" validate:"required"`
	QuestionText          string             `json:"questionText" validate:"required"`
	NextQuestionIfTrueID  *uint              `json:"nextQuestionIfTrueID" validate:"omitempty"`
	NextQuestionIfFalseID *uint              `json:"nextQuestionIfFalseID" validate:"omitempty"`
	CorrectAnswer         string             `json:"correctAnswer" validate:"omitempty"`
	IsDependency          bool               `json:"isDependency" validate:"omitempty"`
	QuestionType          types.QuestionType `json:"questionType" validate:"required,oneof=Conditional ConditionalWithAnswer MultipleChoice MultipleChoiceWithAnswer Descriptive"`
	QuestionOptions       []QuestionOption   `json:"questionOptions" validate:"omitempty"`
}

type QuestionOption struct {
	OptionText     string
	NextQuestionID *uint
}

func TypeToDomainMapper(question types.Question, surveyUUID uuid.UUID) *Question {
	var questions []QuestionOption
	var nextQuestionIfFalseID, nextQuestionIfTrueID *uint
	if question.QuestionType == types.Conditional || question.QuestionType == types.ConditionalWithAnswer {
		nextQuestionIfFalseID = question.Options[1].NextQuestionID
		nextQuestionIfTrueID = question.Options[0].NextQuestionID
	} else {
		for _, option := range question.Options {
			questions = append(questions, QuestionOption{OptionText: option.OptionText, NextQuestionID: option.NextQuestionID})
		}
	}
	return &Question{
		ID:                    question.ID,
		SurveyUUID:            surveyUUID,
		QuestionText:          question.QuestionText,
		NextQuestionIfTrueID:  nextQuestionIfTrueID,
		NextQuestionIfFalseID: nextQuestionIfFalseID,
		IsDependency:          question.IsDependency,
		QuestionType:          question.QuestionType,
		QuestionOptions:       questions,
		CorrectAnswer:         question.CorrectAnswer,
	}
}

func DomainToTypeMapper(question Question, surveyID uint) types.Question {
	var questions []types.QuestionOption
	for _, option := range question.QuestionOptions {
		questions = append(questions, types.QuestionOption{OptionText: option.OptionText, NextQuestionID: option.NextQuestionID})
	}
	if question.QuestionType == types.Conditional || question.QuestionType == types.ConditionalWithAnswer {
		questions = append(questions, types.QuestionOption{OptionText: "true", NextQuestionID: question.NextQuestionIfTrueID})
		questions = append(questions, types.QuestionOption{OptionText: "false", NextQuestionID: question.NextQuestionIfTrueID})
	}
	return types.Question{
		SurveyID:      surveyID,
		QuestionText:  question.QuestionText,
		IsDependency:  question.IsDependency,
		QuestionType:  question.QuestionType,
		Options:       questions,
		CorrectAnswer: question.CorrectAnswer,
	}
}
