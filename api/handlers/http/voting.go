package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/porseOnline/api/service"
	"github.com/porseOnline/pkg/adapters/storage/types"
	"github.com/porseOnline/pkg/logger"
)

func Vote(srv *service.VoteService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req types.Vote
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		err := srv.Vote(c.UserContext(), &req)
		if err != nil {
			logger.Error("can not vote", nil)
			return fiber.ErrInternalServerError
		}
		return nil
	}
}
