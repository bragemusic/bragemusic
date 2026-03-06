-- migrate:up
CREATE TABLE IF NOT EXISTS likes (
    id          TEXT PRIMARY KEY,
    track_id    TEXT NOT NULL,
    owner       TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,

    FOREIGN KEY (track_id)
      REFERENCES tracks(id)
      ON DELETE CASCADE
);

CREATE INDEX idx_likes_track_id ON likes(track_id);

CREATE INDEX idx_likes_owner_created_at
ON likes(owner, created_at DESC);

CREATE UNIQUE INDEX idx_likes_track_owner
ON likes(track_id, owner);

-- migrate:down

