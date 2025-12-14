package adapter

// 1. The Standard Data Format
// No matter where the song comes from, we convert it to this clean format.
type SongMetadata struct {
	ISRC   string
	Title  string
	Artist string
}

// 2. The Contract
// Every platform (Spotify, Deezer, etc.) MUST be able to do these two things.
type PlatformAdapter interface {
	// A. If the user sends a link from this platform, extract the info.
	// Example: "spotify.com/track/123" -> {Title: "Hello", ISRC: "US123..."}
	GetMetadata(url string) (*SongMetadata, error)

	// B. If we have a song info, find the link on this platform.
	// Example: {ISRC: "US123..."} -> "spotify.com/track/123"
	Search(isrc, title, artist string) (string, error)
}
