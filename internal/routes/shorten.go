package routes

import (
	"github.com/ebriussenex/shrt/internal/handlers"
	"github.com/ebriussenex/shrt/internal/shorten"
	"github.com/gofiber/fiber/v2"
)

func ShortenRoutes(a *fiber.App, shortenService *shorten.Service) {
	route := a.Group("/api/v1")
	route.Post("/short", handlers.HandleShorten(shortenService))
	route.Get("/:id", handlers.HandleRedirect(shortenService))
}