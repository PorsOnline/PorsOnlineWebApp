package domain

type Vote struct {
	ID             string
	UserID         string
	SurveyID       string
	QuestionID     string
	TextResponse   string
	SelectedOption string
}
