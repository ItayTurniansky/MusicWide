package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	// Ensure these import paths are correct for your project structure
	"github.com/ItayTurniansky/musicwide/internal/server"
	"github.com/ItayTurniansky/musicwide/internal/service"
)

func main() {
	// 1. Load Environment Variables (.env is used for local development only)
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found or failed to load. Using system environment variables.")
	}

	// 2. Initialize the Music Service
	musicService, err := service.NewMusicService()
	if err != nil {
		log.Fatal("Failed to start Music Service:", err)
	}
	log.Println("Music Service Initialized!")

	// 3. Start the Server
	s := server.NewServer(musicService)

	// Get the PORT from the environment (used by Render), defaulting to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s (Public URL will be provided by Render)", port)

	// Corrected ListenAndServe: This single line starts the server using
	// the custom router (s.Router) and logs a fatal error if it fails.
	log.Fatal(http.ListenAndServe(":"+port, s.Router))
}
