package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

type AppleAdapter struct {
	Client *http.Client
}

func NewAppleAdapter() *AppleAdapter {
	return &AppleAdapter{
		Client: &http.Client{},
	}
}

// 1. GetMetadata: Apple Link -> ISRC
func (a *AppleAdapter) GetMetadata(link string) (*SongMetadata, error) {
	// A. EXTRACT TRACK ID
	var trackID string
	u, err := url.Parse(link)
	if err == nil {
		trackID = u.Query().Get("i")
	}

	// Fallback if "i=" is missing
	if trackID == "" {
		pathIDRegex := regexp.MustCompile(`/([0-9]+)`)
		matches := pathIDRegex.FindAllStringSubmatch(link, -1)
		if len(matches) > 0 {
			trackID = matches[len(matches)-1][1]
		}
	}

	if trackID == "" {
		return nil, fmt.Errorf("could not find apple track ID")
	}

	// B. EXTRACT COUNTRY CODE (Defaults to "us" if not found)
	country := "us"
	countryRegex := regexp.MustCompile(`music\.apple\.com/([a-z]{2})/`)
	countryMatches := countryRegex.FindStringSubmatch(link)
	if len(countryMatches) >= 2 {
		country = countryMatches[1]
	}

	// C. CALL API
	apiURL := fmt.Sprintf("https://itunes.apple.com/lookup?id=%s&country=%s", trackID, country)

	resp, err := a.Client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			ISRC       string `json:"isrc"`
			TrackName  string `json:"trackName"`
			ArtistName string `json:"artistName"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.ResultCount == 0 || len(result.Results) == 0 {
		return nil, fmt.Errorf("song not found")
	}

	track := result.Results[0]

	if track.TrackName == "" {
		return nil, fmt.Errorf("found ID but title was empty")
	}

	return &SongMetadata{
		ISRC:   track.ISRC,
		Title:  track.TrackName,
		Artist: track.ArtistName,
	}, nil
}

// 2. Search: ISRC -> Apple Link
func (a *AppleAdapter) Search(isrc, title, artist string) (string, string, error) {
	// Attempt 1: Search by ISRC
	if isrc != "" {
		link, err := a.doSearch(isrc)
		if err == nil {
			return link, "ISRC (Exact Match)", nil
		}
	}

	// Attempt 2: Text Search
	textQuery := fmt.Sprintf("%s %s", title, artist)
	link, err := a.doSearch(textQuery)
	if err == nil {
		return link, "Text Search", nil
	}

	return "", "", fmt.Errorf("song not found on apple music")
}

func (a *AppleAdapter) doSearch(query string) (string, error) {
	encodedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://itunes.apple.com/search?term=%s&entity=song&limit=1&country=us", encodedQuery)

	resp, err := a.Client.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			TrackViewUrl string `json:"trackViewUrl"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.ResultCount == 0 || len(result.Results) == 0 {
		return "", fmt.Errorf("not found")
	}

	return result.Results[0].TrackViewUrl, nil
}
