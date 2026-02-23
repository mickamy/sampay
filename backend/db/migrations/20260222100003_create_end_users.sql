-- migrate:up
CREATE TABLE end_users
(
    user_id    CHAR(26)    NOT NULL PRIMARY KEY REFERENCES users (id),
    slug       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE end_users;
