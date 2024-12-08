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

func CreateRole(svc *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req types.Role
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := svc.CreateRole(c.UserContext(), *mapper.RoleStorage2Domain(req))
		if err != nil {
			logger.Error("error in creating role", nil)
			return err
		}
		logger.Info("role created successfully", nil)
		return c.JSON(resp)
	}
}

func GetRole(svc *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		resp, err := svc.GetRole(c.UserContext(), domain.RoleID(id))
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			logger.Error("error in fetching role", nil)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		logger.Info("fetched role successfully", nil)
		return c.JSON(resp)
	}
}

func UpdateRole(svc *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req types.Role
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		err := svc.UpdateRole(c.UserContext(), *mapper.RoleStorage2Domain(req))
		if err != nil {
			logger.Error("error in updating role", nil)
			return err
		}
		logger.Info("role updated successfully", nil)
		return nil
	}
}

func DeleteRole(svc *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = svc.DeleteRole(c.UserContext(), domain.RoleID(id))
		if err != nil {
			logger.Error("error in deleting role", nil)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		logger.Info("deleted role successfully", nil)
		return nil
	}
}

func AssignRoleToUser(svc *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId, err := c.ParamsInt("roleId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		userId, err := c.ParamsInt("userId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = svc.AssignRoleToUser(c.UserContext(), domain.RoleID(roleId), domain.UserID(userId))
		if err != nil {
			logger.Error("error in assining role to user", nil)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		logger.Info("assigned role to user successfully", nil)
		return nil
	}
}
