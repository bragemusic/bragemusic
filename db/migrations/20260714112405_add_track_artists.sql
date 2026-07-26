-- migrate:up

CREATE TABLE IF NOT EXISTS track_artists (
    id              TEXT PRIMARY KEY,
    track_id        TEXT NOT NULL,
    artist_id       TEXT NOT NULL,
    role            TEXT NOT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (track_id, artist_id),

    FOREIGN KEY (track_id)
      REFERENCES tracks(id)
      ON DELETE CASCADE,

    FOREIGN KEY (artist_id)
      REFERENCES artists(id)
      ON DELETE RESTRICT
);

CREATE INDEX idx_track_artists_artist
  ON track_artists (artist_id);

-- migrate:down

