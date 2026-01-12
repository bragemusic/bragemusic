-- migrate:up

CREATE TABLE IF NOT EXISTS albums (
    id              TEXT PRIMARY KEY,              -- UUID
    musicbrainz_id  TEXT UNIQUE,
    name            TEXT NOT NULL,
    sort_name       TEXT NOT NULL,
    release_date    DATETIME,
    tracks          INTEGER,
    discs           INTEGER,
    description     TEXT,
    owner           TEXT NOT NULL,
    public          BOOLEAN,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down

