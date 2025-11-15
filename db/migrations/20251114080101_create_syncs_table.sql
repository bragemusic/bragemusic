-- migrate:up

CREATE TABLE IF NOT EXISTS syncs (
    id              TEXT PRIMARY KEY,              -- UUID
    artists_created INTEGER,
    artists_updated INTEGER,
    albums_created  INTEGER,
    albums_updated  INTEGER,
    tracks_created  INTEGER,
    tracks_updated  INTEGER,
    synced_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down

