CREATE TABLE IF NOT EXISTS user_domains (
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    registration_date DATE DEFAULT NOW()::DATE,
    version BIGINT DEFAULT 1
);
