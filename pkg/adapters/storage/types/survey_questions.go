package types

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	SurveyID              uint
	QuestionText          string
	Order                 int
	NextQuestionIfTrueID  *uint
	NextQuestionIfFalseID *uint
	CorrectAnswer         string
	QuestionType          QuestionType
	Options               []QuestionOption
	IsDependency          bool
}

type QuestionOption struct {
	gorm.Model
	QuestionID uint
	OptionText string
	IsCorrect  bool
}

type QuestionType string

const (
	Conditional              QuestionType = "Conditional"
	ConditionalWithAnswer    QuestionType = "ConditionalWithAnswer"
	MultipleChoice           QuestionType = "MultipleChoice"
	MultipleChoiceWithAnswer QuestionType = "MultipleChoiceWithAnswer"
	Descriptive              QuestionType = "Descriptive"
)
