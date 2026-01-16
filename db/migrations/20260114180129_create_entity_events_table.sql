-- migrate:up
CREATE TABLE entity_events (
    id UUID NOT NULL,
    event_type TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    event_time TIMESTAMP NOT NULL,
    PRIMARY KEY (entity_type, id)
);

-- migrate:down

