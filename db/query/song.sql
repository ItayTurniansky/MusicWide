-- name: CreateSong :one
INSERT INTO songs (
  isrc, title, artist, spotify_url, apple_url, deezer_url, youtube_url
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSongByISRC :one
SELECT * FROM songs
WHERE isrc = $1 LIMIT 1;