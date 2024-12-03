package http

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/porseOnline/api/service"
	"github.com/porseOnline/pkg/adapters/storage/types"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func SendMessage(srv *service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req types.Notification
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		req.ID = uuid.Must(uuid.NewRandom()).String()
		req.Read = false
		req.Create_at = time.Now()
		err := srv.SendMessage(c.UserContext(), &req)
		if err != nil {
			log.Error("can not send message")
			return err
		}
		return nil
	}

}

func GetUnreadMessages(srv *service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("user_id")
		resp, err := srv.GetUnreadMessages(c.UserContext(), userID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return c.JSON(resp)

	}
}
