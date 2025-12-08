-- migrate:up

CREATE TABLE play_history (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    track_id TEXT NOT NULL,
    played_at DATETIME NOT NULL
);

CREATE INDEX idx_play_history_user_track
    ON play_history (user_id, track_id);

CREATE INDEX idx_play_history_user_id
    ON play_history (user_id);

CREATE INDEX idx_play_history_track_id
    ON play_history (track_id);

CREATE INDEX idx_play_history_played_at
    ON play_history (played_at);

-- migrate:down

