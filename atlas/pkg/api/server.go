package api

import (
	"atlas/pkg/mongodb"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	addr string
	cfg  fiber.Config
	app  *fiber.App
	mc   *mongodb.MongoClient
}

func NewServer(addr string, mc *mongodb.MongoClient, cfg fiber.Config) *Server {
	return &Server{
		addr: addr,
		cfg:  cfg,
		app:  fiber.New(cfg),
		mc:   mc,
	}
}

func (s *Server) Listen() {
	s.app.Listen(s.addr)
}
