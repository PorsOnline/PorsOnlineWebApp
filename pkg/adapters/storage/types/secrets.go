package types

import "gorm.io/gorm"

type Secrets struct {
	gorm.Model
	UserID   string `json:"userID"  yaml:"user_id"`
	ServeyID string `json:"serveyID"  yaml:"servey_id"`
	Secret   string `json:"secret"  yaml:"secret"`
}
