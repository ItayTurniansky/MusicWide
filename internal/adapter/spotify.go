package adapter

import (
	"context"
	"fmt"
	"os"
	"regexp" // <--- Added Regex package

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyAdapter struct {
	Client *spotify.Client
}

func NewSpotifyAdapter() (*SpotifyAdapter, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("spotify credentials missing in .env")
	}

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	httpClient := config.Client(context.Background())
	client := spotify.New(httpClient)

	return &SpotifyAdapter{Client: client}, nil
}

// 1. Compile the Regex once (Global variable)
// This looks for "track/" or "track:" followed by 22 letters/numbers
var spotifyTrackRegex = regexp.MustCompile(`(?:track[:/])([a-zA-Z0-9]{22})`)

func (s *SpotifyAdapter) GetMetadata(url string) (*SongMetadata, error) {
	// 2. Use Regex to find the ID
	matches := spotifyTrackRegex.FindStringSubmatch(url)

	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid spotify track link: could not find 22-char ID")
	}

	trackID := matches[1] // The part inside the ([...])

	// 3. Ask Spotify for details
	track, err := s.Client.GetTrack(context.Background(), spotify.ID(trackID))
	if err != nil {
		return nil, err
	}

	isrc := track.ExternalIDs["isrc"]
	if isrc == "" {
		return nil, fmt.Errorf("no ISRC found for this spotify track")
	}

	return &SongMetadata{
		ISRC:   isrc,
		Title:  track.Name,
		Artist: track.Artists[0].Name,
	}, nil
}

// Search takes an ISRC and finds the Spotify Link
// Returns: (Link, Method, Error)
func (s *SpotifyAdapter) Search(isrc, title, artist string) (string, string, error) {
	// 1. Search by ISRC (ONLY if we have a valid ISRC)
	if isrc != "" {
		query := "isrc:" + isrc
		results, err := s.Client.Search(context.Background(), query, spotify.SearchTypeTrack)
		if err == nil && results.Tracks != nil && len(results.Tracks.Tracks) > 0 {
			return results.Tracks.Tracks[0].ExternalURLs["spotify"], "ISRC (Exact Match)", nil
		}
	}

	// 2. Fallback: Search by Title + Artist
	query := fmt.Sprintf("track:%s artist:%s", title, artist)
	results, err := s.Client.Search(context.Background(), query, spotify.SearchTypeTrack)
	if err != nil {
		return "", "", err
	}

	if results.Tracks != nil && len(results.Tracks.Tracks) > 0 {
		return results.Tracks.Tracks[0].ExternalURLs["spotify"], "Text Search", nil
	}

	return "", "", fmt.Errorf("song not found on spotify")
}
