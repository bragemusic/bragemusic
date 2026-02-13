-- migrate:up
CREATE TABLE ratings (
    id          TEXT PRIMARY KEY,
    track_id    TEXT NOT NULL,
    rating      INTEGER NOT NULL,
    owner       TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,

    FOREIGN KEY (track_id)
      REFERENCES tracks(id)
      ON DELETE CASCADE
);

CREATE INDEX idx_ratings_track_id ON ratings(track_id);

CREATE INDEX idx_ratings_owner_created_at
ON ratings(owner, created_at DESC);

CREATE UNIQUE INDEX idx_ratings_track_owner
ON ratings(track_id, owner);


-- migrate:down

