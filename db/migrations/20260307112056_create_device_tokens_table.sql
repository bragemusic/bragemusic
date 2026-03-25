-- migrate:up
CREATE TABLE IF NOT EXISTS device_tokens (
    device_id   TEXT NOT NULL,
    token_id    TEXT NOT NULL,
    created_at  DATETIME NOT NULL,

    PRIMARY KEY (device_id, token_id),

    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (token_id) REFERENCES tokens(id) ON DELETE CASCADE
);

-- migrate:down

