-- migrate:up

ALTER TABLE playlists
ADD COLUMN type TEXT NOT NULL DEFAULT 'standard';

CREATE TABLE IF NOT EXISTS smart_playlist_contents (
    id              TEXT PRIMARY KEY,
    playlist_id     TEXT NOT NULL,
    bpm_upper       INTEGER,
    bpm_lower       INTEGER,
    mood_happy      REAL,
    mood_sad        REAL,
    mood_aggressive REAL,
    mood_calm       REAL,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL,

    FOREIGN KEY (playlist_id)
      REFERENCES playlists(id)
      ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS smart_playlist_artists (
    id             TEXT PRIMARY KEY,
    parent_id      TEXT NOT NULL,
    artist_id      TEXT NOT NULL,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL,

    FOREIGN KEY (parent_id)
      REFERENCES smart_playlist_contents(id)
      ON DELETE CASCADE
);

-- migrate:down

