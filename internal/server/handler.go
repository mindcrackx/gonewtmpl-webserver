package server

import (
	"log/slog"
	"net/http"
)

func (s *Server) HandleHelloGet(w http.ResponseWriter, _ *http.Request) {
	const tmplName = "hello"

	if err := s.templates.Render(w, tmplName, nil); err != nil {
		slog.Error("HandleHelloGet template render", "template", tmplName, "err", err)
		return
	}
}
