package domain

import "github.com/porseOnline/pkg/adapters/storage/types"

type Question struct {
	ID                    uint
	SurveyID              uint
	QuestionText          string
	NextQuestionIfTrueID  *uint
	NextQuestionIfFalseID *uint
	IsDependency          *bool
	QuestionType          types.QuestionType
	QuestionOptions       []QuestionOption
}

type QuestionOption struct {
	OptionText string
	IsCorrect  bool
}

func TypeToDomainMapper(question types.Question) *Question {
	var questions []QuestionOption
	for _, option := range question.Options {
		questions = append(questions, QuestionOption{OptionText: option.OptionText, IsCorrect: option.IsCorrect})
	}
	return &Question{
		ID:                    question.ID,
		SurveyID:              question.SurveyID,
		QuestionText:          question.QuestionText,
		NextQuestionIfTrueID:  question.NextQuestionIfTrueID,
		NextQuestionIfFalseID: question.NextQuestionIfFalseID,
		IsDependency:          question.IsDependency,
		QuestionType:          question.QuestionType,
		QuestionOptions:       questions,
	}
}

func DomainToTypeMapper(question Question) types.Question {
	var questions []types.QuestionOption
	for _, option := range question.QuestionOptions {
		questions = append(questions, types.QuestionOption{OptionText: option.OptionText, IsCorrect: option.IsCorrect})
	}
	return types.Question{
		SurveyID:              question.SurveyID,
		QuestionText:          question.QuestionText,
		NextQuestionIfTrueID:  question.NextQuestionIfTrueID,
		NextQuestionIfFalseID: question.NextQuestionIfFalseID,
		IsDependency:          question.IsDependency,
		QuestionType:          question.QuestionType,
		Options:               questions,
	}
}
