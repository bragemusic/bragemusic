-- migrate:up

CREATE TABLE IF NOT EXISTS imports (
    id              TEXT PRIMARY KEY,
    musicbrainz_id  TEXT,
    owner           TEXT NOT NULL,
    filename        TEXT NOT NULL,
    type            TEXT NOT NULL,
    state           TEXT NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down

