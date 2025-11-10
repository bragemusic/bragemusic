-- migrate:up

CREATE TABLE IF NOT EXISTS artists (
    id              TEXT PRIMARY KEY,              -- UUID
    musicbrainz_id  TEXT,
    name            TEXT NOT NULL,
    sort_name       TEXT NOT NULL,
    country         TEXT,
    year_started    INTEGER,
    year_ended      INTEGER,
    description     TEXT,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS artists_updated_at_trigger
AFTER UPDATE ON artists
FOR EACH ROW
BEGIN
    UPDATE artists SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;


-- migrate:down

