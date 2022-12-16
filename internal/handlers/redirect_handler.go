package handlers

import (
	"context"
	"errors"
	"log"

	"github.com/ebriussenex/shrt/internal/model"
	"github.com/gofiber/fiber/v2"
)

type redirecter interface {
	Redirect(ctx context.Context, id string) (string, error)
}

func HandleRedirect(redirecter redirecter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		redirectUrl, err := redirecter.Redirect(c.UserContext(), id);
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"err": true,
					"msg": err.Error(),
				}) 
			}
			log.Printf("error getting redirect url for %q: %v", id, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": true,
				"msg": "server error during redirection",
			}) 
		}
		log.Printf("redirecting to %s", redirectUrl)
		return c.Redirect(redirectUrl, fiber.StatusMovedPermanently)
		
	}
}