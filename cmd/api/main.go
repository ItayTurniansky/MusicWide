package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/ItayTurniansky/musicwide/internal/server"
	"github.com/ItayTurniansky/musicwide/internal/service"
)

func main() {
	// 1. Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Initialize the Music Service
	musicService, err := service.NewMusicService()
	if err != nil {
		log.Fatal(" Failed to start Music Service:", err)
	}
	fmt.Println(" Music Service Initialized!")

	// 3. Start the Server
	s := server.NewServer(musicService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
	fmt.Printf(" MusicWide Server running on http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}
