package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

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

func UpdateQuestion(svc *service.QuestionService) fiber.Handler {
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
		response, err := svc.UpdateQuestion(c.UserContext(), &req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(response)
	}
}

func DeleteQuestion(svc *service.QuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		param := c.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			return fiber.ErrBadRequest
		}
		err = svc.DeleteQuestion(c.UserContext(), uint(id))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON("deleted successfully")
	}
}

func GetQuestion(svc *service.QuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), "UserID", c.Locals("UserID"))
		var req domain.UserQuestionStep
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		resp, err := svc.GetNextQuestion(ctx, req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(resp)
	}
}
