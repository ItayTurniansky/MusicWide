package server

import (
	"encoding/json"
	"net/http"

	db "github.com/ItayTurniansky/musicwide/db/sqlc"
)

// Define the JSON structure we expect from the user
type CreateSongRequest struct {
	ISRC    string `json:"isrc"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Spotify string `json:"spotify"`
	Apple   string `json:"apple"`
	Deezer  string `json:"deezer"`
	Youtube string `json:"youtube"`
}

// The Logic Function
func (s *Server) CreateSong(w http.ResponseWriter, r *http.Request) {
	var req CreateSongRequest

	// 1. Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Prepare data for the database
	arg := db.CreateSongParams{
		Isrc:       req.ISRC,
		Title:      req.Title,
		Artist:     req.Artist,
		SpotifyUrl: req.Spotify,
		AppleUrl:   req.Apple,
		DeezerUrl:  req.Deezer,
		YoutubeUrl: req.Youtube,
	}

	// 3. Save to Database
	song, err := s.db.CreateSong(r.Context(), arg)
	if err != nil {
		http.Error(w, "Failed to save song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Respond with the saved song
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}
