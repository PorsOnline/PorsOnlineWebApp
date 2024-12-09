package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/porseOnline/api/service"
	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/pkg/adapters/storage/mapper"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"github.com/porseOnline/pkg/logger"
)

type UserPermissionValidationRequest struct {
	Resource string `json:"resource"`
	Scope    string `json:"scope"`
	Group    string `json:"group"`
}

func CreatePermission(svc *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req types.Permission
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := svc.CreatePermission(c.UserContext(), *mapper.PermissionStorage2Domain(req))
		if err != nil {
			logger.Error("error in creating permission", nil)
			return err
		}
		logger.Info("permission created successfully", nil)
		return c.JSON(resp)
	}
}

func GetUserPermissions(svc *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		resp, err := svc.GetUserPermissions(c.UserContext(), domain.UserID(id))
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			logger.Error("error in fetching user permissions", nil)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		logger.Info("fetched user permissions successfully", nil)
		return c.JSON(resp)
	}
}

func GetPermissionByID(svc *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		resp, err := svc.GetPermissionByID(c.UserContext(), domain.PermissionID(id))
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			logger.Error("error in fetching permission", nil)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		logger.Info("fetched permission successfully", nil)
		return c.JSON(resp)
	}
}

func UpdatePermission(svc *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req types.Permission
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		err := svc.UpdatePermission(c.UserContext(), *mapper.PermissionStorage2Domain(req))
		if err != nil {
			logger.Error("error in updating permission", nil)
			return err
		}
		logger.Info("permission updated successfully", nil)
		return nil
	}
}

func DeletePermission(svc *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = svc.DeletePermission(c.UserContext(), domain.PermissionID(id))
		if err != nil {
			logger.Error("error in deleting permission", nil)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		logger.Info("deleted permission successfully", nil)
		return nil
	}
}

func ValidateUserPermission(svc *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req UserPermissionValidationRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		userId, err := c.ParamsInt("userId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		valid, err := svc.ValidateUserPermission(c.UserContext(), domain.UserID(userId), req.Resource, req.Scope, req.Group)
		if err != nil {
			logger.Error("error in validating user permission", nil)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if !valid {
			logger.Error("user do not have access to this resource", nil)
			return fiber.NewError(fiber.StatusNotAcceptable, "Forbidden")
		}
		logger.Info("validate user permission successfully", nil)
		return nil
	}
}

func AssignPermissionToUser(svc *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req []domain.PermissionDetails
		if err := c.BodyParser(&req); err != nil {
			logger.Error("error in parse assign permission body", nil)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err := svc.AssignPermissionToUser(c.UserContext(), req)
		if err != nil {
			logger.Error("error in assining permission to user", nil)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		logger.Info("assigned permission to user successfully", nil)
		return nil
	}
}
