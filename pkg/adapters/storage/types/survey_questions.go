package types

import (
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	SurveyID      uint
	QuestionText  string
	Order         int
	CorrectAnswer string
	QuestionType  QuestionType
	Options       []QuestionOption
	IsDependency  bool
}

type QuestionOption struct {
	gorm.Model
	QuestionID     uint
	OptionText     string
	NextQuestionID *uint
}

type UserQuestionStep struct {
	gorm.Model
	SurveyID   uint
	UserID     uint
	QuestionID uint
	Action     Action
}

type QuestionType string

const (
	ConditionalMultipleChoice QuestionType = "ConditionalMultipleChoice"
	MultipleChoice            QuestionType = "MultipleChoice"
	MultipleChoiceWithAnswer  QuestionType = "MultipleChoiceWithAnswer"
	Descriptive               QuestionType = "Descriptive"
)

type Action string

const (
	Forward  Action = "forward"
	Backward Action = "backward"
)
