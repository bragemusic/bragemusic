-- migrate:up

CREATE TABLE IF NOT EXISTS album_tracks (
    id              TEXT NOT NULL,
    album_id        TEXT NOT NULL,
    track_id        TEXT NOT NULL,
    disc_number     INTEGER NOT NULL DEFAULT 1,
    track_number    INTEGER NOT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (album_id, track_id),

    FOREIGN KEY (album_id)
      REFERENCES albums(id)
      ON DELETE CASCADE,

    FOREIGN KEY (track_id)
      REFERENCES tracks(id)
      ON DELETE RESTRICT
);

CREATE INDEX idx_album_tracks_album
  ON album_tracks (album_id);

-- migrate:down

