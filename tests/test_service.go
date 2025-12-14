package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ItayTurniansky/musicwide/internal/service"
	"github.com/joho/godotenv"
)

func test_service() {
	// 1. Load Environment Variables (.env)
	// We need this for the Spotify and YouTube API keys.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Start the "Brain" (MusicService)
	// This function connects to Spotify and YouTube internally.
	svc, err := service.NewMusicService()
	if err != nil {
		log.Fatal("❌ Failed to start service:", err)
	}
	fmt.Println("✅ MusicService Started Successfully")

	// ---------------------------------------------------------
	// 3. The Test Input
	// You can change this link to ANY platform (Spotify, Apple, etc.)
	// Let's use the YouTube link that was causing trouble earlier.
	// ---------------------------------------------------------
	inputLink := "https://www.youtube.com/watch?v=Sy4-JH449LY"

	fmt.Printf("\n🚀 Processing Link: %s\n", inputLink)

	// 4. Run the Conversion
	// This single line does ALL the work:
	// - Detects it's YouTube
	// - Extracts "Scapegoat (Remix)"
	// - Searches Spotify, Deezer, and Apple
	result, err := svc.ConvertLink(inputLink)
	if err != nil {
		log.Fatal("❌ Conversion Failed:", err)
	}

	// 5. Print the Result as JSON
	// This makes it look like the data your future website will receive.
	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonBytes))
}
