-- migrate:up

CREATE TABLE IF NOT EXISTS tracks (
    id              TEXT PRIMARY KEY,              -- UUID
    title           TEXT NOT NULL,
    album_id        TEXT NOT NULL,
    musicbrainz_id  TEXT,
    track_artist    TEXT,
    track_number    INTEGER,
    disc_number     INTEGER,
    genre           TEXT,
    year            INTEGER,
    composer        TEXT,
    comment         TEXT,
    duration_ms     INTEGER,
    bitrate         INTEGER,
    sample_rate     INTEGER,
    file_path       TEXT NOT NULL,
    file_size       INTEGER,
    mime_type       TEXT,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS tracks_updated_at_trigger
AFTER UPDATE ON tracks
FOR EACH ROW
BEGIN
    UPDATE tracks SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- migrate:down

