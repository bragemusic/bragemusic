-- migrate:up

CREATE TABLE IF NOT EXISTS album_artists (
    id              TEXT PRIMARY KEY,
    album_id        TEXT NOT NULL,
    artist_id       TEXT NOT NULL,
    role            TEXT NOT NULL,
    position        INTEGER NOT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (album_id, artist_id, role),

    FOREIGN KEY (album_id)
      REFERENCES albums(id)
      ON DELETE CASCADE,

    FOREIGN KEY (artist_id)
      REFERENCES artists(id)
      ON DELETE RESTRICT
);

CREATE INDEX idx_album_artists_album
  ON album_artists (album_id);

-- migrate:down

