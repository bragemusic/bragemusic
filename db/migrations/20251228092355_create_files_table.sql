-- migrate:up

CREATE TABLE IF NOT EXISTS media_files (
    id              TEXT PRIMARY KEY,              -- UUID
    duration_ms     INTEGER NOT NULL,
    bitrate         INTEGER NOT NULL,
    sample_rate     INTEGER NOT NULL,
    file_size       INTEGER NOT NULL,
    codec           TEXT NOT NULL,
    checksum        TEXT NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down

