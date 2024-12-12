package http

import (
	"errors"
	"net/http"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/porseOnline/api/service"
	"github.com/porseOnline/internal/question"
	"github.com/porseOnline/internal/question/domain"
	"github.com/porseOnline/internal/survey"
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
			if errors.Is(err, &question.ErrBadRequest{}) || errors.Is(err, &question.ErrNoQuestionFound{}){
				fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
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
		surveyParam := c.Params("surveyID")
		id, err := strconv.Atoi(param)
		surveyID, err := strconv.Atoi(surveyParam)
		if err != nil {
			return fiber.ErrBadRequest
		}
		err = svc.DeleteQuestion(c.UserContext(), uint(id), uint(surveyID))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON("deleted successfully")
	}
}

func GetNextQuestion(svc *service.QuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req domain.UserQuestionStep
		surveyParam := c.Params("surveyID")
		surveyID, err := strconv.Atoi(surveyParam)
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		userID, err := strconv.Atoi(c.Locals("UserID").(string))
		resp, err := svc.GetNextQuestion(c.UserContext(), req, uint(userID), uint(surveyID))
		if err != nil {
			if errors.Is(err, &survey.ErrBadRequest{}) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())	
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(resp)
	}
}
