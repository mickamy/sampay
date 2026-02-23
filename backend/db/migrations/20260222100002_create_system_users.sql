-- migrate:up
CREATE TABLE system_users
(
    user_id CHAR(26) NOT NULL PRIMARY KEY REFERENCES users (id),
    name    TEXT     NOT NULL
);

CREATE UNIQUE INDEX idx_system_users_name ON system_users (name);

-- migrate:down
DROP TABLE system_users;
