package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// Import your database package
	db "github.com/ItayTurniansky/musicwide/db/sqlc"
)

type Server struct {
	Router *chi.Mux
	db     *db.Queries
}

func NewServer(database *db.Queries) *Server {
	s := &Server{
		Router: chi.NewRouter(),
		db:     database,
	}

	s.mountMiddleware()
	s.mountRoutes()

	return s
}

func (s *Server) mountMiddleware() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.AllowContentType("application/json"))
}

func (s *Server) mountRoutes() {
	s.Router.Get("/health", s.HealthHandler)

	// This function (s.CreateSong) is defined in song_handler.go
	// Go finds it automatically because they are in the same package.
	s.Router.Post("/songs", s.CreateSong)
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"MusicWide is alive!"}`))
}
