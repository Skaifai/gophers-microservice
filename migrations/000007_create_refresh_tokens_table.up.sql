CREATE TABLE IF NOT EXISTS refresh_tokens (
    key TEXT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    token_string TEXT NOT NULL
);
