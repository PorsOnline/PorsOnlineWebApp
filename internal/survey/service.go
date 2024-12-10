package survey

import (
	"context"
	"errors"

	"github.com/porseOnline/internal/survey/domain"
	surveyPort "github.com/porseOnline/internal/survey/port"
	permissionDomain "github.com/porseOnline/internal/user/domain"
	userPort "github.com/porseOnline/internal/user/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type service struct {
	repo              surveyPort.Repo
	permissionService userPort.PermissionService
}

type ErrBadRequest struct{}

func (m *ErrBadRequest) Error() string {
	return "survey not found"
}

func NewService(repo surveyPort.Repo, permissionService userPort.PermissionService) surveyPort.Service {
	return &service{repo: repo, permissionService: permissionService}
}

func (ss *service) CreateSurvey(ctx context.Context, survey domain.Survey, userID uint) (*domain.Survey, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	survey.UUID = newUUID
	typeSurvey := domain.DomainToTypeMapper(survey)
	typeSurvey.UserID = userID
	createdSurvey, err := ss.repo.Create(ctx, typeSurvey, survey.TargetCities)
	if err != nil {
		return nil, err
	}
	ss.permissionService.AssignSurveyPermissionsToOwner(ctx, surveyPermissions(), userID, createdSurvey.ID)
	return domain.TypeToDomainMapper(*createdSurvey), nil
}

func (ss *service) UpdateSurvey(ctx context.Context, survey domain.Survey, id uint) (*domain.Survey, error) {
	typeSurvey := domain.DomainToTypeMapper(survey)
	typeSurvey.ID = id
	updatedSurvey, err := ss.repo.Update(ctx, typeSurvey, survey.TargetCities)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &ErrBadRequest{}
		}
		return nil, err
	}
	return domain.TypeToDomainMapper(*updatedSurvey), nil
}

func (ss *service) GetAllSurveys(ctx context.Context, page, pageSize int) ([]domain.Survey, error) {
	surveys, err := ss.repo.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	var domainSurveys []domain.Survey
	for _, survey := range surveys {
		domainSurveys = append(domainSurveys, *domain.TypeToDomainMapper(survey))
	}
	return domainSurveys, nil
}

func (ss *service) GetSurveyByUUID(ctx context.Context, surveyUUID uuid.UUID) (*domain.Survey, error) {
	survey, err := ss.repo.GetByUUID(ctx, surveyUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &ErrBadRequest{}
		}
		return nil, err
	}
	return domain.TypeToDomainMapper(*survey), nil
}

func (ss *service) CancelSurvey(ctx context.Context, id uint) error {
	err := ss.repo.Cancel(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &ErrBadRequest{}
		}
	}
	return nil
}

func (ss *service) DeleteSurvey(ctx context.Context, id uint) error {
	err := ss.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &ErrBadRequest{}
		}
	}
	return nil
}

func (ss *service) GetSurveyByID(ctx context.Context, surveyID uint) (*domain.Survey, error) {
	survey, err := ss.repo.GetByID(ctx, surveyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &ErrBadRequest{}
		}
		return nil, err
	}
	return domain.TypeToDomainMapper(*survey), nil
}

func surveyPermissions() []permissionDomain.Permission {
	return []permissionDomain.Permission{
		{Resource: "/api/v1/survey/:uuid", Scope: "read"},
		{Resource: "/api/v1/survey", Scope: "update"},
		{Resource: "/api/v1/survey/cancel/:uuid", Scope: "create"},
		{Resource: "/api/v1/survey/:uuid", Scope: "delete"},
		{Resource: "/api/v1/survey", Scope: "read"},
		{Resource: "/api/v1/survey/:id/question", Scope: "create"},
		{Resource: "/api/v1/survey/:id/question/:id", Scope: "delete"},
		{Resource: "/api/v1/survey/:id/question", Scope: "update"},
		{Resource: "/api/v1/survey/:id/question/get-next", Scope: "read"},
	}
}
