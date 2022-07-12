DROP TABLE IF NOT EXISTS users, logins, passwords;

CREATE TABLE IF NOT EXISTS passwords (
    id BIGSERIAL,
    hash text NOT NULL,
    salt text NOT NULL,
    is_active bool NOT NULL,
    created BIGINT DEFAULT extract(epoch from now()),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS logins (
    id BIGSERIAL,
    login text NOT NULL,
    password_id BIGINT,
    PRIMARY KEY(id),
    FOREIGN KEY(password_id) REFERENCES passwords(id)
);

CREATE TABLE IF NOT EXISTS Users (
    id BIGSERIAL,
    is_disabled bool,
    password_id BIGINT,
    login_id BIGINT,
    PRIMARY KEY(id),
    FOREIGN KEY(login_id) REFERENCES logins(id),
    FOREIGN KEY(password_id) REFERENCES passwords(id)
);
