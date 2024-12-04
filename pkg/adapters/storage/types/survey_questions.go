package types

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	SurveyID              uint
	QuestionText          string
	Order                 int
	NextQuestionIfTrueID  *uint
	NextQuestionIfFalseID *uint
	QuestionType          QuestionType
	Options               []QuestionOption
	IsDependency          *bool
}

type QuestionOption struct {
	gorm.Model
	QuestionID uint
	OptionText string
	IsCorrect  bool
}

type QuestionType string

const (
	Conditional           QuestionType = "conditional"
	ConditionalWithAnswer QuestionType = "conditional_with_answer"
	Optional              QuestionType = "optional"
	OptionalWithAnswer    QuestionType = "optional_with_answer"
	Descriptive           QuestionType = "descriptive"
)
