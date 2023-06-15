CREATE TABLE IF NOT EXISTS user_auths (
    domain_user_id BIGINT NOT NULL PRIMARY KEY,
    role TEXT NOT NULL DEFAULT 'USER',
    password TEXT NOT NULL,
    activation_link TEXT UNIQUE,
    activated BOOLEAN DEFAULT false,
    CONSTRAINT fk_user_profiles
        FOREIGN KEY (domain_user_id)
        REFERENCES user_domains(id)
        ON DELETE CASCADE
);
