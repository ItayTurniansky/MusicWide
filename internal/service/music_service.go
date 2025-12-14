package service

import (
	"fmt"
	"sync" // We will use this to make searching faster!

	"github.com/ItayTurniansky/musicwide/internal/adapter"
)

type MusicService struct {
	Spotify *adapter.SpotifyAdapter
	Deezer  *adapter.DeezerAdapter
	Apple   *adapter.AppleAdapter
	Youtube *adapter.YoutubeAdapter
}

func NewMusicService() (*MusicService, error) {
	spot, err := adapter.NewSpotifyAdapter()
	if err != nil {
		return nil, fmt.Errorf("spotify init failed: %w", err)
	}

	yt, err := adapter.NewYoutubeAdapter()
	if err != nil {
		return nil, fmt.Errorf("youtube init failed: %w", err)
	}

	return &MusicService{
		Spotify: spot,
		Deezer:  adapter.NewDeezerAdapter(),
		Apple:   adapter.NewAppleAdapter(),
		Youtube: yt,
	}, nil
}

// ConvertResult is the final output
type ConvertResult struct {
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Spotify string `json:"spotify"`
	Deezer  string `json:"deezer"`
	Apple   string `json:"apple"`
	Youtube string `json:"youtube"`
}

// ConvertLink does everything in one go: Extract -> Search All
func (s *MusicService) ConvertLink(inputLink string) (*ConvertResult, error) {
	// 1. IDENTIFY & EXTRACT
	platform := adapter.IdentifyPlatform(inputLink)
	fmt.Printf("   🔍 Detected Platform: %s\n", platform)

	var meta *adapter.SongMetadata
	var err error

	switch platform {
	case adapter.PlatformSpotify:
		meta, err = s.Spotify.GetMetadata(inputLink)
	case adapter.PlatformDeezer:
		meta, err = s.Deezer.GetMetadata(inputLink)
	case adapter.PlatformApple:
		meta, err = s.Apple.GetMetadata(inputLink)
	case adapter.PlatformYoutube:
		meta, err = s.Youtube.GetMetadata(inputLink)
	default:
		return nil, fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	// 2. SEARCH EVERYWHERE (Using Go Routines for Speed)
	// Since we are doing 4 network calls, doing them in parallel makes the API much faster.
	result := &ConvertResult{
		Title:  meta.Title,
		Artist: meta.Artist,
	}

	var wg sync.WaitGroup
	wg.Add(4)

	// Search Spotify
	go func() {
		defer wg.Done()
		link, _, _ := s.Spotify.Search(meta.ISRC, meta.Title, meta.Artist)
		result.Spotify = link
	}()

	// Search Deezer
	go func() {
		defer wg.Done()
		link, _, _ := s.Deezer.Search(meta.ISRC, meta.Title, meta.Artist)
		result.Deezer = link
	}()

	// Search Apple
	go func() {
		defer wg.Done()
		link, _, _ := s.Apple.Search(meta.ISRC, meta.Title, meta.Artist)
		result.Apple = link
	}()

	// Search YouTube
	go func() {
		defer wg.Done()
		link, _, _ := s.Youtube.Search(meta.ISRC, meta.Title, meta.Artist)
		result.Youtube = link
	}()

	wg.Wait() // Wait for all searches to finish

	return result, nil
}
