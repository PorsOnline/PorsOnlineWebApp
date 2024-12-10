package storage

import (
	"context"

	surveyPort "github.com/porseOnline/internal/survey/port"
	"github.com/porseOnline/pkg/adapters/storage/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type surveyRepo struct {
	db *gorm.DB
}

func NewSurveyRepo(db *gorm.DB) surveyPort.Repo {
	return &surveyRepo{db: db}
}

func (sr *surveyRepo) Delete(ctx context.Context, id uint) error {
	return	sr.db.Where("id = ?", id).Delete(&types.Survey{Model: gorm.Model{ID: id}}).Error
}

func (sr *surveyRepo) Cancel(ctx context.Context, id uint) error {
	var survey types.Survey
	err := sr.db.Model(&types.Survey{}).Where("id = ?", id).First(&survey).Error
	if err != nil {
		return err
	}
	survey.IsActive = false
	return sr.db.Model(&types.Survey{}).Where("id = ?", id).Save(&survey).Error
}

func (sr *surveyRepo) GetByUUID(ctx context.Context, uuid uuid.UUID) (*types.Survey, error) {
	var survey *types.Survey
	err := sr.db.Model(&types.Survey{}).Preload("TargetCities").Where("UUID = ?", uuid).First(&survey).Error
	if err != nil {
		return nil, err
	}
	return survey, nil
}

func (sr *surveyRepo) GetAll(ctx context.Context, page, pageSize int) ([]types.Survey, error) {
	var surveys []types.Survey
	err := sr.db.Model(&types.Survey{}).Preload("TargetCities").Limit(pageSize).Offset((page-1)*pageSize).Where("deleted_at is null").Find(&surveys).Error
	if err != nil {
		return nil, err
	}
	return surveys, nil
}

func (sr *surveyRepo) Create(ctx context.Context, survey types.Survey, cities []string) (*types.Survey, error) {
	tx := sr.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	err := tx.Debug().Model(&types.Survey{}).Create(&survey).Error
	if err != nil {
		return nil, err
	}
	for _, city := range cities {
		var typeCity types.SurveyCity
		typeCity.Name = city
		typeCity.SurveyID = survey.ID
		err := tx.Model(&types.SurveyCity{}).Debug().Create(&typeCity).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &survey, nil
}

func (sr *surveyRepo) Update(ctx context.Context, survey types.Survey, cities []string) (*types.Survey, error) {
	var oldSurvey types.Survey
	err := sr.db.Model(&types.Survey{}).Where("id = ?", survey.ID).First(&oldSurvey).Error
	if err != nil {
		return nil, err
	}
	survey.ID = oldSurvey.ID
	tx := sr.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	err = tx.Debug().Model(&types.Survey{}).Where("id = ?", oldSurvey.ID).Save(&survey).Error
	if err != nil {
		return nil, err
	}
	err = sr.db.Model(&types.SurveyCity{}).Where("survey_id = ? ", survey.ID).Delete(&types.SurveyCity{}).Error
	if err != nil {
		tx.Rollback()
		return &survey, err
	}
	for _, city := range cities {
		err = tx.Model(&types.SurveyCity{}).Debug().Create(&types.SurveyCity{SurveyID: survey.ID, Name: city}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &survey, nil
}

func (sr *surveyRepo) GetByID(ctx context.Context, id uint) (*types.Survey, error) {
	var survey *types.Survey
	err := sr.db.Model(&types.Survey{}).Preload("TargetCities").Where("id = ?", id).First(&survey).Error
	if err != nil {
		return nil, err
	}
	return survey, nil
}