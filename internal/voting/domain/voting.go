package domain

import "time"

type Vote struct {
	ID             uint
	UserID         string
	SurveyID       string
	QuestionID     string
	TextResponse   string
	SelectedOption string
	CreatedAt      time.Time
	DeletedAt      time.Time
	UpdatedAt      time.Time
}
