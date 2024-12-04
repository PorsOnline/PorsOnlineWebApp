package http

import (
	"errors"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/porseOnline/api/service"
	"github.com/porseOnline/internal/question/domain"
)

func CreateQuestion(svc *service.QuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req domain.Question
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		validate := validator.New()
		err := validate.Struct(req)
		if err != nil {
			var errs validator.ValidationErrors
			errors.As(err, &errs)
			for _, validationError := range errs {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": map[string]string{
						"field":   validationError.Field(),
						"message": validationError.Error(),
					},
				})
			}
		}
		response, err := svc.CreateQuestion(c.UserContext(), &req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(response)
	}
}