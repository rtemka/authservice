DROP TABLE IF NOT EXISTS users, passwords;

CREATE TABLE IF NOT EXISTS users (
    login text,
    created_at BIGINT DEFAULT extract(epoch from now()),
    is_disabled bool,
    PRIMARY KEY(login)
);

CREATE TABLE IF NOT EXISTS passwords (
    hash text NOT NULL,
    user_login text NOT NULL,
    is_active bool NOT NULL,
    generated_at BIGINT DEFAULT extract(epoch from now()),
    PRIMARY KEY(hash),
    FOREIGN KEY(user_login) REFERENCES users(login)
);