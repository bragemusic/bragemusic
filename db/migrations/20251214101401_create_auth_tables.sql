-- migrate:up

CREATE TABLE IF NOT EXISTS users (
    id           TEXT PRIMARY KEY,        -- UUID
    email        TEXT,                    -- nullable (OAuth)
    username     TEXT,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_email
ON users(email)
WHERE email IS NOT NULL;

-----------
CREATE TABLE IF NOT EXISTS local_credentials (
    user_id        TEXT PRIMARY KEY,
    password_hash  TEXT NOT NULL,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-----------
CREATE TABLE IF NOT EXISTS auth_identities (
    id                TEXT PRIMARY KEY,           -- UUID
    user_id           TEXT NOT NULL,
    provider          TEXT NOT NULL,        -- "local", "github", "authelia"
    provider_user_id  TEXT NOT NULL,
    email             TEXT,
    created_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    UNIQUE (provider, provider_user_id)
);


-----------
CREATE TABLE IF NOT EXISTS tokens (
    id            TEXT PRIMARY KEY,           -- UUID
    user_id       TEXT NOT NULL,
    name          TEXT,                     -- "laptop", "ci-runner", "backup-daemon"
    hash          TEXT NOT NULL,            -- SHA-256(token)
    scopes        TEXT NOT NULL,          -- space-separated scopes
    expires_at    DATETIME,               -- NULL = long-lived
    last_used_at  DATETIME,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_tokens_hash
ON tokens(hash);

-----------
CREATE TABLE IF NOT EXISTS user_scopes (
    user_id     TEXT NOT NULL,
    scope       TEXT NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, scope),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-----------
CREATE TABLE IF NOT EXISTS token_events (
    id          TEXT PRIMARY KEY,           -- UUID
    token_id    TEXT NOT NULL,
    event       TEXT NOT NULL,           -- "created", "revoked", "used"
    ip_address  TEXT,
    user_agent  TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (token_id)
        REFERENCES tokens(id)
        ON DELETE CASCADE
);





-- migrate:down

