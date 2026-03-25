-- migrate:up
CREATE TABLE IF NOT EXISTS devices (
    id                 TEXT PRIMARY KEY,
    name               TEXT NOT NULL,
    type               TEXT NOT NULL,
    interface          TEXT NOT NULL,
    icon               TEXT NOT NULL,
    user_id            TEXT NOT NULL,
    supports_playback  BOOLEAN NOT NULL,
    platform           TEXT NOT NULL,
    version            TEXT NOT NULL,
    last_ip            TEXT NOT NULL,
    last_seen          DATETIME NOT NULL,
    created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- migrate:down
