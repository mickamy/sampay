-- migrate:up
CREATE TABLE users
(
    id         CHAR(26)    NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE users;
