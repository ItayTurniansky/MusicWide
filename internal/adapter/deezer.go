package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings" // <--- Added strings package
)

type DeezerAdapter struct {
	Client *http.Client
}

func NewDeezerAdapter() *DeezerAdapter {
	return &DeezerAdapter{
		Client: &http.Client{},
	}
}

// Regex to find Deezer Track IDs (e.g., deezer.com/track/12345)
var deezerTrackRegex = regexp.MustCompile(`track/(\d+)`)

// 1. GetMetadata: Deezer Link -> ISRC
func (d *DeezerAdapter) GetMetadata(link string) (*SongMetadata, error) {

	// --- FIX START: Handle Short Links ---
	// If the link doesn't look like a standard "track/123" link,
	// we follow the redirect to get the real URL.
	if !strings.Contains(link, "track/") {
		resp, err := d.Client.Get(link) // Get follows redirects automatically
		if err != nil {
			return nil, fmt.Errorf("failed to resolve short link: %w", err)
		}
		// The final URL (after redirects) is here:
		link = resp.Request.URL.String()
		resp.Body.Close() // We don't need the body, just the URL
	}
	// --- FIX END ---

	// A. Find the ID in the URL (Now that we have the full URL)
	matches := deezerTrackRegex.FindStringSubmatch(link)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid deezer link: could not find track ID")
	}
	trackID := matches[1]

	// B. Call Deezer API
	apiURL := fmt.Sprintf("https://api.deezer.com/track/%s", trackID)
	resp, err := d.Client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("deezer api returned status %d", resp.StatusCode)
	}

	// C. Parse the Response
	var result struct {
		ISRC   string `json:"isrc"`
		Title  string `json:"title"`
		Artist struct {
			Name string `json:"name"`
		} `json:"artist"`
		Error struct {
			Type string `json:"type"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error.Type != "" {
		return nil, fmt.Errorf("deezer api error: %s", result.Error.Type)
	}

	return &SongMetadata{
		ISRC:   result.ISRC,
		Title:  result.Title,
		Artist: result.Artist.Name,
	}, nil
}

// 2. Search: ISRC -> Deezer Link
func (d *DeezerAdapter) Search(isrc, title, artist string) (string, string, error) {
	// Attempt 1: Search by ISRC (Only if valid)
	if isrc != "" {
		query := fmt.Sprintf("isrc:%s", isrc)
		link, err := d.doSearch(query)
		if err == nil {
			return link, "ISRC (Exact Match)", nil
		}
	}

	// Attempt 2: Strict Text Search
	strictQuery := fmt.Sprintf("artist:\"%s\" track:\"%s\"", artist, title)
	link, err := d.doSearch(strictQuery)
	if err == nil {
		return link, "Strict Name Search", nil
	}

	// Attempt 3: Fuzzy Search
	fuzzyQuery := fmt.Sprintf("%s %s", artist, title)
	link, err = d.doSearch(fuzzyQuery)
	if err == nil {
		return link, "Fuzzy Name Search", nil
	}

	return "", "", fmt.Errorf("no song found on deezer")
}

// Helper function
func (d *DeezerAdapter) doSearch(query string) (string, error) {
	encodedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://api.deezer.com/search?q=%s", encodedQuery)

	resp, err := d.Client.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Link string `json:"link"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Data) == 0 {
		return "", fmt.Errorf("no song found")
	}

	return result.Data[0].Link, nil
}
