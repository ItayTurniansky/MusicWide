CREATE TABLE "songs" (
  "id" bigserial PRIMARY KEY,
  "isrc" varchar UNIQUE NOT NULL, -- ISRC is the best unique ID for a song
  "title" varchar NOT NULL,
  "artist" varchar NOT NULL,      -- Added Artist because a song name alone isn't enough
  "spotify_url" varchar NOT NULL,
  "apple_url" varchar NOT NULL,
  "deezer_url" varchar NOT NULL,
  "youtube_url" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Index for searching by ISRC (fast lookups)
CREATE INDEX ON "songs" ("isrc");
-- Index for searching by Title (if people search by name)
CREATE INDEX ON "songs" ("title");