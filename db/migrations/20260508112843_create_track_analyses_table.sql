-- migrate:up
CREATE TABLE IF NOT EXISTS track_analyses (
    id             TEXT PRIMARY KEY,
    bpm            INTEGER NOT NULL,
    key            TEXT NOT NULL,
    key_scale      TEXT NOT NULL,
    key_confidence REAL NOT NULL,
    loudness       INTEGER NOT NULL,
    energy         REAL NOT NULL,
    danceability   REAL NOT NULL,
    mood_happy     REAL NOT NULL,
    mood_sad       REAL NOT NULL,
    mood_aggresive REAL NOT NULL,
    mood_calm      REAL NOT NULL,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL,

    FOREIGN KEY (id)
      REFERENCES tracks(id)
      ON DELETE CASCADE
);


-- migrate:down

