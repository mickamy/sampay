-- migrate:up
CREATE TABLE line_friendships
(
    end_user_id CHAR(26)    NOT NULL PRIMARY KEY REFERENCES end_users (user_id),
    is_friend   BOOLEAN     NOT NULL DEFAULT FALSE,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS line_friendships;
