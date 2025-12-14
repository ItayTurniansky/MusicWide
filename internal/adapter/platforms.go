package adapter

import "strings"

// Define a type for our platforms (like an Enum)
type Platform string

const (
	PlatformSpotify Platform = "spotify"
	PlatformDeezer  Platform = "deezer"
	PlatformApple   Platform = "apple"
	PlatformYoutube Platform = "youtube"
	PlatformUnknown Platform = "unknown"
)

// Helper function to figure out which platform a link belongs to
func IdentifyPlatform(url string) Platform {
	lowerURL := strings.ToLower(url)

	if strings.Contains(lowerURL, "spotify.com") {
		return PlatformSpotify
	}
	if strings.Contains(lowerURL, "deezer.com") {
		return PlatformDeezer
	}
	if strings.Contains(lowerURL, "apple.com") || strings.Contains(lowerURL, "music.apple.com") {
		return PlatformApple
	}
	if strings.Contains(lowerURL, "youtube.com") || strings.Contains(lowerURL, "youtu.be") {
		return PlatformYoutube
	}

	return PlatformUnknown
}
