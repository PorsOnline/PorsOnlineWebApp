package http

import (
	"strconv"

	"github.com/porseOnline/api/service"
	"github.com/porseOnline/internal/survey/domain"

	"github.com/gofiber/fiber/v2"
)

func CreateSurvey(svcGetter ServiceGetter[*service.SurveyService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		svc := svcGetter(c.UserContext())
		var req domain.Survey
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		userID, err := strconv.Atoi(c.Locals("UserID").(string))
		response, err := svc.CreateSurvey(c.UserContext(), &req, uint(userID))
		if err != nil {
			// if errors.Is(err, service.ErrUserCreationValidation) {
			// 	return fiber.NewError(fiber.StatusBadRequest, err.Error())
			// }

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(response)
	}
}

func GetSurvey(svcGetter ServiceGetter[*service.SurveyService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var param = c.Params("surveyID")
		surveyID, err := strconv.Atoi(param)

		svc := svcGetter(c.UserContext())

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		resp, err := svc.GetSurvey(c.Context(), uint(surveyID))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(resp)
	}
}

func UpdateSurvey(svcGetter ServiceGetter[*service.SurveyService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var param = c.Params("surveyID")
		surveyID, err := strconv.Atoi(param)

		svc := svcGetter(c.UserContext())

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		var req domain.Survey
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		response, err := svc.UpdateSurvey(c.UserContext(), &req, uint(surveyID))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(response)
	}
}

func CancelSurvey(svcGetter ServiceGetter[*service.SurveyService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var param = c.Params("surveyID")
		surveyID, err := strconv.Atoi(param)

		svc := svcGetter(c.UserContext())

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = svc.CancelSurvey(c.Context(), uint(surveyID))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON("successful")
	}
}

func DeleteSurvey(svcGetter ServiceGetter[*service.SurveyService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var param = c.Params("surveyID")
		surveyID, err := strconv.Atoi(param)

		svc := svcGetter(c.UserContext())

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = svc.DeleteSurvey(c.Context(), uint(surveyID))
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

func GetAllSurveys(svcGetter ServiceGetter[*service.SurveyService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
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
