-- migrate:up
CREATE TABLE oauth_accounts
(
    id          CHAR(26)    NOT NULL PRIMARY KEY,
    end_user_id CHAR(26)    NOT NULL REFERENCES end_users (user_id),
    provider    TEXT        NOT NULL,
    uid         TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (provider, uid)
);
CREATE INDEX idx_oauth_accounts_end_user_id ON oauth_accounts (end_user_id);

-- migrate:down
DROP TABLE oauth_accounts;
