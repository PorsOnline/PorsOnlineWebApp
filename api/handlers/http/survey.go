package http

import (
	"github.com/porseOnline/api/service"
	"github.com/porseOnline/internal/survey/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


func CreateSurvey(svc *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req domain.Survey
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		response, err := svc.CreateSurvey(c.UserContext(), &req)
		if err != nil {
			// if errors.Is(err, service.ErrUserCreationValidation) {
			// 	return fiber.NewError(fiber.StatusBadRequest, err.Error())
			// }

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(response)
	}
}

func GetSurvey(svc *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var param = c.Params("uuid")
		uuid, err := uuid.Parse(param)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		resp, err := svc.GetSurvey(c.Context(), uuid)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(resp)
	}
}

func UpdateSurvey(svc *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var param = c.Params("uuid")
		uuid, err := uuid.Parse(param)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		var req domain.Survey
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		req.UUID = uuid
		response, err := svc.UpdateSurvey(c.UserContext(), &req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(response)
	}
}

func CancelSurvey(svc *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var param = c.Params("uuid")
		uuid, err := uuid.Parse(param)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = svc.CancelSurvey(c.Context(), uuid)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON("successful")
	}
}

func DeleteSurvey(svc *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var param = c.Params("uuid")
		uuid, err := uuid.Parse(param)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = svc.DeleteSurvey(c.Context(), uuid)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON("successful")
	}
}

type PaginationQuery struct {
    Page int `query:"page" default:"1" validate:"gt=0"`
    Size int `query:"size" default:"10" validate:"gt=0"`
    // SortBy string `query:"sortBy" default:"name" validate:"oneof=id name country"`
}

func GetAllSurveys(svc *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var paginationQuery PaginationQuery
		err := c.QueryParser(&paginationQuery)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		resp, err := svc.GetAllSurveys(c.Context(), paginationQuery.Page, paginationQuery.Size)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(resp)
	}
}