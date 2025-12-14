package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ItayTurniansky/musicwide/internal/service"
)

type Server struct {
	Router *chi.Mux
	Music  *service.MusicService
}

// NewServer now only needs the MusicService
func NewServer(music *service.MusicService) *Server {
	s := &Server{
		Router: chi.NewRouter(),
		Music:  music,
	}

	s.mountMiddleware()
	s.mountRoutes()

	return s
}

func (s *Server) mountMiddleware() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	// We removed AllowContentType to allow simple browser GET requests
}

func (s *Server) mountRoutes() {
	s.Router.Get("/health", s.HealthHandler)
	s.Router.Get("/convert", s.HandleConvert)
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"MusicWide is alive!"}`))
}

// HandleConvert accepts a link and returns all matches
func (s *Server) HandleConvert(w http.ResponseWriter, r *http.Request) {
	// CORS Headers (Crucial for Frontend)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	link := r.URL.Query().Get("link")
	if link == "" {
		http.Error(w, `{"error": "Missing 'link' parameter"}`, http.StatusBadRequest)
		return
	}

	fmt.Printf(" Processing Request: %s\n", link)

	// Call Service
	result, err := s.Music.ConvertLink(link)
	if err != nil {
		// Log the error to console so you can debug
		fmt.Printf(" Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Return Success
	json.NewEncoder(w).Encode(result)
}
