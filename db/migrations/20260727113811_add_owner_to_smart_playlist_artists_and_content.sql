-- migrate:up
ALTER TABLE smart_playlist_contents
ADD COLUMN owner TEXT NOT NULL DEFAULT "";

ALTER TABLE smart_playlist_artists
ADD COLUMN owner TEXT NOT NULL DEFAULT "";

-- migrate:down

