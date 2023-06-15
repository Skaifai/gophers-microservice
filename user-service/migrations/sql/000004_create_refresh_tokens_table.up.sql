CREATE TABLE IF NOT EXISTS refresh_tokens (
   key TEXT NOT NULL PRIMARY KEY,
   created_at TIMESTAMP NOT NULL,
   token_string TEXT NOT NULL
);
