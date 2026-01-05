-- migrate:up

CREATE TABLE IF NOT EXISTS sync_items (
    id              TEXT PRIMARY KEY,
    sync_id         TEXT NOT NULL,
    item_id         TEXT NOT NULL,
    type            TEXT NOT NULL,
    state           TEXT NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
);

-- migrate:down

