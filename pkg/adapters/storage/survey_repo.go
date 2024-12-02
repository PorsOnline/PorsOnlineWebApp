package storage

import (
	surveyPort "PorsOnlineWebApp/internal/survey/port"
	"PorsOnlineWebApp/pkg/adapters/storage/types"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type surveyRepo struct {
	db *gorm.DB
}

func NewSurveyRepo(db *gorm.DB) surveyPort.Repo {
	return &surveyRepo{db: db}
}

func (sr *surveyRepo) Delete(ctx context.Context, uuid uuid.UUID) error {
	return	sr.db.Delete(&types.Survey{UUID: uuid}).Error
}

func (sr *surveyRepo) Cancel(ctx context.Context, uuid uuid.UUID) error {
	var survey types.Survey
	err := sr.db.Model(&types.Survey{}).Where("UUID = ?", uuid).Find(&survey).Error
	if err != nil {
		return err
	}
	survey.IsActive = false
	return sr.db.Model(&types.Survey{}).Where("UUID = ?", uuid).Save(&survey).Error
}

func (sr *surveyRepo) Get(ctx context.Context, uuid uuid.UUID) (*types.Survey, error) {
	var survey types.Survey
	err := sr.db.Model(&types.Survey{}).Where("UUID = ?", uuid).Find(&survey).Error
	return &survey, err
}

func (sr *surveyRepo) GetAll(ctx context.Context, page, pageSize int) ([]types.Survey, error) {
	var surveys []types.Survey
	err := sr.db.Model(&types.Survey{}).Limit(pageSize).Offset((page-1)*pageSize).Find(&surveys).Error
	if err != nil {
		return nil, err
	}
	return surveys, nil
}

func (sr *surveyRepo) Create(ctx context.Context, survey types.Survey) (*types.Survey, error) {
	return &survey, sr.db.Model(&types.Survey{}).Create(survey).Error
}

func (sr *surveyRepo) Update(ctx context.Context, survey types.Survey) (*types.Survey, error) {
	var oldSurvey types.Survey
	err := sr.db.Model(&types.Survey{}).Where("UUID = ?", survey.UUID).First(&oldSurvey).Error
	if err != nil {
		return &survey, err
	}
	return &survey, sr.db.Model(&types.Survey{}).Save(survey).Error	
}