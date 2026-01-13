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

-- migrate:down
