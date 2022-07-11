DROP TABLE IF NOT EXISTS logins, passwords;

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
