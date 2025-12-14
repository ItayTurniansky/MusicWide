package main

import (
	"fmt"
	"log"

	"github.com/ItayTurniansky/musicwide/internal/adapter"
	"github.com/joho/godotenv"
)

func test_adpaters() {
	// 1. Setup
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize All Adapters
	spotify, err := adapter.NewSpotifyAdapter()
	if err != nil {
		log.Fatal("❌ Spotify Error:", err)
	}

	deezer := adapter.NewDeezerAdapter()
	apple := adapter.NewAppleAdapter()

	youtube, err := adapter.NewYoutubeAdapter()
	if err != nil {
		log.Fatal("❌ YouTube Error:", err)
	}

	fmt.Println("✅ SYSTEM ONLINE: All 4 Adapters Connected")
	fmt.Println("==================================================")

	// ---------------------------------------------------------
	// SCENARIO 1: SPOTIFY INPUT (Despacito)
	// ---------------------------------------------------------
	fmt.Println("\n🔵 SCENARIO 1: Input from SPOTIFY")
	sInput := "https://open.spotify.com/track/0ZBjWfNaRotYYTGTZ2lEtA"
	fmt.Printf("   Input: %s\n", sInput)

	meta, err := spotify.GetMetadata(sInput)
	if err != nil {
		fmt.Printf("   ❌ Extract Failed: %v\n", err)
	} else {
		printMetadata(meta)
		// Convert to others
		testConversion("Deezer ", meta, deezer)
		testConversion("Apple  ", meta, apple)
		testConversion("YouTube", meta, youtube)
	}

	// ---------------------------------------------------------
	// SCENARIO 2: DEEZER INPUT (Shape of You)
	// ---------------------------------------------------------
	fmt.Println("\n\n🟣 SCENARIO 2: Input from DEEZER")
	dInput := "https://www.deezer.com/track/140397657"
	fmt.Printf("   Input: %s\n", dInput)

	meta, err = deezer.GetMetadata(dInput)
	if err != nil {
		fmt.Printf("   ❌ Extract Failed: %v\n", err)
	} else {
		printMetadata(meta)
		// Convert to others
		testConversion("Spotify", meta, spotify)
		testConversion("Apple  ", meta, apple)
		testConversion("YouTube", meta, youtube)
	}

	// ---------------------------------------------------------
	// SCENARIO 3: APPLE INPUT (Blinding Lights)
	// ---------------------------------------------------------
	fmt.Println("\n\n🔴 SCENARIO 3: Input from APPLE MUSIC")
	// Using the US link (which your adapter now handles correctly)
	aInput := "https://music.apple.com/us/album/scapegoat-yebisu303-remix/1196952243?i=1196952288&uo=4"
	fmt.Printf("   Input: %s\n", aInput)

	meta, err = apple.GetMetadata(aInput)
	if err != nil {
		fmt.Printf("   ❌ Extract Failed: %v\n", err)
	} else {
		printMetadata(meta)
		// Convert to others
		testConversion("Spotify", meta, spotify)
		testConversion("Deezer ", meta, deezer)
		testConversion("YouTube", meta, youtube)
	}

	// ---------------------------------------------------------
	// SCENARIO 4: YOUTUBE INPUT (Uptown Funk)
	// ---------------------------------------------------------
	fmt.Println("\n\n🟡 SCENARIO 4: Input from YOUTUBE")
	// "Mark Ronson - Uptown Funk (Official Video) ft. Bruno Mars"
	// This tests if your Title Cleaner works!
	yInput := "https://www.youtube.com/watch?v=Sy4-JH449LY"
	fmt.Printf("   Input: %s\n", yInput)

	meta, err = youtube.GetMetadata(yInput)
	if err != nil {
		fmt.Printf("   ❌ Extract Failed: %v\n", err)
	} else {
		printMetadata(meta)
		// Convert to others (All must use Text Search since YouTube has no ISRC)
		testConversion("Spotify", meta, spotify)
		testConversion("Deezer ", meta, deezer)
		testConversion("Apple  ", meta, apple)
	}

	fmt.Println("\n==================================================")
	fmt.Println("✅ ALL TESTS COMPLETED")
}

// Interface to make the test function generic
type Searcher interface {
	Search(isrc, title, artist string) (string, string, error)
}

func testConversion(platformName string, meta *adapter.SongMetadata, adapter Searcher) {
	link, method, err := adapter.Search(meta.ISRC, meta.Title, meta.Artist)
	if err != nil {
		fmt.Printf("   ❌ %s: Failed (%v)\n", platformName, err)
	} else {
		fmt.Printf("   ✅ %s: Found via %s\n", platformName, method)
		fmt.Printf("      Link: %s\n", link)
	}
}

func printMetadata(m *adapter.SongMetadata) {
	fmt.Printf("   ℹ️  Extracted: \"%s\" by %s\n", m.Title, m.Artist)
	if m.ISRC != "" {
		fmt.Printf("      ISRC: %s\n", m.ISRC)
	} else {
		fmt.Println("      ISRC: [Not Found] (Will rely on text search)")
	}
}
