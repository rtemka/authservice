DROP TABLE IF EXISTS users, passwords;

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

INSERT INTO users(login, created_at, is_disabled)
VALUES ('login1', 1658141437, FALSE),
('login2', 1658141437, FALSE);

INSERT INTO passwords(user_login, hash, generated_at, is_active)
VALUES ('login1', 'h1', 1658141437, TRUE),
('login2', 'h2', 1658141437, TRUE); 