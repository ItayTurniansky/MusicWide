package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type YoutubeAdapter struct {
	ApiKey string
	Client *http.Client
}

func NewYoutubeAdapter() (*YoutubeAdapter, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY is missing in .env")
	}
	return &YoutubeAdapter{
		ApiKey: apiKey,
		Client: &http.Client{},
	}, nil
}

// Regex to find YouTube Video IDs
// Handles:
// - music.youtube.com/watch?v=ID
// - youtube.com/watch?v=ID
// - youtu.be/ID
var ytVideoRegex = regexp.MustCompile(`(?:v=|be/)([\w-]{11})`)

// 1. GetMetadata: YouTube Link -> Title/Artist
func (y *YoutubeAdapter) GetMetadata(link string) (*SongMetadata, error) {
	matches := ytVideoRegex.FindStringSubmatch(link)
	if len(matches) < 2 {
		return nil, fmt.Errorf("could not find youtube video ID")
	}
	videoID := matches[1]

	// Call YouTube API
	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s",
		videoID, y.ApiKey,
	)

	resp, err := y.Client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			Snippet struct {
				Title        string `json:"title"`
				ChannelTitle string `json:"channelTitle"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	}

	// Clean the Title (Remove "Official Video", etc.)
	rawTitle := result.Items[0].Snippet.Title
	cleanTitle := cleanYoutubeTitle(rawTitle)
	channel := result.Items[0].Snippet.ChannelTitle

	// Channel names often have " - Topic" appended for auto-generated music
	channel = strings.Replace(channel, " - Topic", "", 1)

	return &SongMetadata{
		ISRC:   "",
		Title:  cleanTitle,
		Artist: channel,
	}, nil
}

// 2. Search: Artist + Title -> YouTube Link
func (y *YoutubeAdapter) Search(isrc, title, artist string) (string, string, error) {
	// We search for "Artist - Title Audio" to prioritize the song
	searchQuery := fmt.Sprintf("%s - %s Audio", artist, title)

	encodedQuery := url.QueryEscape(searchQuery)

	// API Call
	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&type=video&maxResults=1&key=%s",
		encodedQuery, y.ApiKey,
	)

	resp, err := y.Client.Get(apiURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			Id struct {
				VideoId string `json:"videoId"`
			} `json:"id"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if len(result.Items) == 0 {
		return "", "", fmt.Errorf("no video found")
	}

	videoID := result.Items[0].Id.VideoId

	// CHANGE: Back to standard YouTube link
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID), "Text Search", nil
}

// Replace the old cleanYoutubeTitle function with this one:
func cleanYoutubeTitle(title string) string {
	// List of specific junk phrases to remove (case-insensitive)
	junkPhrases := []string{
		`\(Official Video\)`, `\[Official Video\]`,
		`\(Official Audio\)`, `\[Official Audio\]`,
		`\(Video\)`, `\[Video\]`,
		`\(Lyrics\)`, `\[Lyrics\]`,
		`\(HQ\)`, `\[HQ\]`,
		`\(4K\)`, `\[4K\]`,
	}

	clean := title
	for _, phrase := range junkPhrases {
		re := regexp.MustCompile(`(?i)` + phrase) // (?i) makes it case-insensitive
		clean = re.ReplaceAllString(clean, "")
	}

	// Remove extra whitespace that might be left behind
	return strings.TrimSpace(clean)
}
