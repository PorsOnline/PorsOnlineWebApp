package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/porseOnline/api/service"
	"github.com/porseOnline/pkg/adapters/storage/types"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func SendMessage(svcGetter ServiceGetter[*service.NotificationService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		srv := svcGetter(c.UserContext())
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

func GetUnreadMessages(svcGetter ServiceGetter[*service.NotificationService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		srv := svcGetter(c.UserContext())
		userID := c.Params("user_id")
		resp, err := srv.GetUnreadMessages(c.UserContext(), userID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return c.JSON(resp)

	}
}
