-- migrate:up

CREATE TABLE IF NOT EXISTS albums (
    id              TEXT PRIMARY KEY,              -- UUID
    musicbrainz_id  TEXT UNIQUE,
    name            TEXT NOT NULL,
    sort_name       TEXT NOT NULL,
    artist_id       TEXT NOT NULL,
    release_date    DATETIME,
    tracks          INTEGER,
    discs           INTEGER,
    description     TEXT,
    owner           TEXT NOT NULL,
    public          BOOLEAN,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS albums_updated_at_trigger
AFTER UPDATE ON albums
FOR EACH ROW
BEGIN
    UPDATE albums SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- migrate:down

