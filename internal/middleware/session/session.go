package sessionmiddleware

import (
	"lms-backend/internal/api"
	"lms-backend/internal/session"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SessionMiddleware(c *fiber.Ctx) error {
	// skip auth routes - /api/v1/auth/*
	paths := strings.Split(c.Path(), "/")
	if len(paths) >= 1 && paths[1] == "swagger" {
		return c.Next()
	}

	if len(paths) >= 3 && paths[3] == "auth" || paths[3] == "health" {
		return c.Next()
	}

	sess, err := session.Store.Get(c)
	if err != nil {
		return err
	}

	token := sess.Get(session.CookieKey)
	if token == nil {
		err := c.JSON(api.Response{
			Messages: []api.Message{api.InfoMessage("User is not logged in")},
		})
		if err != nil {
			return err
		}
		if err := sess.Destroy(); err != nil {
			return err
		}
		return fiber.NewError(fiber.StatusUnauthorized, "User is not logged in")
	}

	return c.Next()
}
