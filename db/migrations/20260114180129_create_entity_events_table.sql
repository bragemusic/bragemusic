-- migrate:up
CREATE TABLE entity_events (
    id TEXT PRIMARY KEY,
    item_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    event_time TIMESTAMP NOT NULL
);

-- migrate:down

