package http

import (
	"strconv"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/porseOnline/api/service"
	"github.com/porseOnline/internal/user/domain"
	"github.com/porseOnline/pkg/jwt"
)

func TraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("X-Trace-ID", traceID)

		c.Locals("traceID", traceID)

		return c.Next()
	}
}

func newAuthMiddleware(secret []byte) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: secret},
		Claims:      &jwt.UserClaims{},
		TokenLookup: "header:Authorization",
		SuccessHandler: func(ctx *fiber.Ctx) error {
			userClaims := userClaims(ctx)
			if userClaims == nil {
				return fiber.ErrUnauthorized
			}
			ctx.Locals("UserID", strconv.Itoa(int(userClaims.UserID)))
			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		},
		AuthScheme: "Bearer",
	})
}

func PermissionMiddleware(permissionService *service.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := strconv.Atoi(c.Locals("UserID").(string))
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		surveyID := c.Params("surveyID")

		valid, err := permissionService.ValidateUserPermission(c.UserContext(), domain.UserID(userID), c.Path(), method2ScopeMapper(c.Route().Method), "", surveyID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if !valid {
			return fiber.NewError(fiber.StatusForbidden, "Permission Denied")
		}
		return c.Next()
	}
}

func method2ScopeMapper(method string) string {
	if method == "GET" {
		return "read"
	}

	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "HEAD":
		return "create"
	case "DELETE":
		return "delete"
	case "PUT":
		return "update"
	case "PATCH":
		return "patch"
	default:
		return ""
	}
}
