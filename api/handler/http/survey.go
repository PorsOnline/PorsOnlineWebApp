package http

import (
	"PorsOnlineWebApp/api/service"
	"PorsOnlineWebApp/internal/survey/domain"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
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
		responsBody, err := json.Marshal(response)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(responsBody)
	}
}