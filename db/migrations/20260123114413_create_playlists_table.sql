-- migrate:up
CREATE TABLE IF NOT EXISTS playlists (
    id              TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT,
    owner           TEXT NOT NULL,
    public          BOOLEAN,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS playlist_tracks (
    id              TEXT PRIMARY KEY,
    playlist_id     TEXT NOT NULL,
    album_track_id  TEXT NOT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (album_track_id)
      REFERENCES album_tracks(id)
      ON DELETE CASCADE
);

CREATE INDEX idx_playlist_tracks_playlist
  ON playlist_tracks (playlist_id);


-- migrate:down

