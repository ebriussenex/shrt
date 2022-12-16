package handlers

import (
	"context"
	"errors"
	"log"
	"strings"

	config "github.com/ebriussenex/shrt/internal/configs"
	"github.com/ebriussenex/shrt/internal/model"
	"github.com/ebriussenex/shrt/internal/shorten"
	"github.com/ebriussenex/shrt/internal/validator"
	"github.com/gofiber/fiber/v2"
)

type shortener interface {
	Shorten(context.Context, model.ShortenedReq) (*model.Shortened, error)
}

type ShortenRequest struct {
	Url string `json:"url" validate:"required,url"`
	Id  string `json:"id,omitempty" validate:"omitempty,alphanum"`
}

// type shortenResponse struct {
// 	Url string `json:"url,omitempty"`
// 	Msg string `json:"msg,omitempty"`
// }

func HandleShorten(shortener shortener) fiber.Handler {
	v := validator.NewValidator()
	return func(c *fiber.Ctx) error {
		var req ShortenRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"err": true,
				"msg": validator.ValidatorErrors(err),
			})
		}

		if err := v.Validate(req); err != nil {
			return err
		}

		// userToken, ok := c.Get("user").(*jwt.Token)
		// if !ok {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 		"err": true,
		// 		"msg": "failed to get user from context",
		// 	})
		// }

		// userClaims, ok := userToken.Claims.(*model.UserClaims)

		var id string
		if strings.TrimSpace(req.Id) != "" {
			id = req.Id
		}

		input := model.ShortenedReq{
			Url:       req.Url,
			Id:        id,
			CreatedBy: "", //userClaims.User.GitHubLogin,
		}

		shortening, err := shortener.Shorten(c.UserContext(), input)
		if err != nil {
			if errors.Is(err, model.ErrIdAlreadyExists) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"err": true,
					"msg": err.Error(),
				}) 
			}

			log.Printf("shortener.Shorten error during link shortening: %v, input: %v", err, input)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": true,
				"msg": "server error during link shortening",
			}) 
		}

		shortUrl, err := shorten.CreateResUrl(config.Get().BaseUrl, shortening.Id)
		if err != nil {
			log.Printf("shorten.CreateResUrl error during creating res url: %v", err)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": true,
				"msg": "server error during creating result link",
			}) 
		}

		return c.JSON(fiber.Map{
			"error": false,
			"status": "success",
			"msg":   nil,
			"shorten":  shortUrl,
		})
	}
}
