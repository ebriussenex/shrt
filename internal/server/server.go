package server

import (
	"github.com/ebriussenex/shrt/internal/routes"
	"github.com/ebriussenex/shrt/internal/shorten"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	f *fiber.App
	shortService *shorten.Service
}

func New(shortService *shorten.Service) *Server {
	s := &Server{
		shortService: shortService,
	}
	s.setupFiber()
	return s
}

func (s *Server) setupFiber() {
	s.f = fiber.New()
	routes.ShortenRoutes(s.f, s.shortService)
}