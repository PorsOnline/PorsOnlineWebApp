package types

import "gorm.io/gorm"

type Vote struct {
	gorm.Model
	UserID         string `json:"userID"  yaml:"userID"`
	SurveyID       string `json:"surveyID"  yaml:"surveyID"`
	QuestionID     string `json:"questionID"  yaml:"questionID"`
	TextResponse   string `json:"textResponse"  yaml:"textResponse"`
	SelectedOption string `json:"selectedOption"  yaml:"selectedOption"`
}
