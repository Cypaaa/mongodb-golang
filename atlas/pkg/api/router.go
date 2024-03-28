package api

import "atlas/pkg/api/handler"

// setup the routes
func (s *Server) Setup() {
	routerPoll(s)
	routerQuestion(s)
	routerChoice(s)
	routerAnswer(s)
}

func routerPoll(s *Server) {
	handler := handler.NewPollHandler(s.mc, s.mc.Database("poll").Collection("poll"))
	s.app.Get("/poll/:id", handler.Get)
	s.app.Post("/poll", handler.Create)
	s.app.Put("/poll", handler.Update)
	s.app.Delete("/poll", handler.Delete)
}

func routerQuestion(s *Server) {
	handler := handler.NewQuestionHandler(s.mc, s.mc.Database("poll").Collection("question"))
	s.app.Post("/question", handler.Create)
	s.app.Put("/question", handler.Update)
	s.app.Delete("/question", handler.Delete)
	s.app.Get("/question/:id", handler.Get)
}

func routerChoice(s *Server) {
	handler := handler.NewChoiceHandler(s.mc, s.mc.Database("poll").Collection("choice"))
	s.app.Post("/choice", handler.Create)
	s.app.Put("/choice", handler.Update)
	s.app.Delete("/choice", handler.Delete)
	s.app.Get("/choice/:id", handler.Get)
}

func routerAnswer(s *Server) {
	handler := handler.NewAnswerHandler(s.mc, s.mc.Database("poll").Collection("answer"))
	s.app.Post("/answer", handler.Create)
	s.app.Delete("/answer", handler.Delete)
	s.app.Get("/answer/:id", handler.Get)
}
