-- migrate:up
CREATE TABLE events
(
    id           CHAR(26)    NOT NULL PRIMARY KEY,
    user_id      CHAR(26)    NOT NULL REFERENCES end_users (user_id),
    title        TEXT        NOT NULL,
    description  TEXT        NOT NULL DEFAULT '',
    total_amount INT         NOT NULL,
    tier_count   INT         NOT NULL DEFAULT 1,
    remainder    INT         NOT NULL DEFAULT 0,
    held_at      TIMESTAMPTZ NOT NULL,
    archived_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_user_id ON events (user_id);

-- migrate:down
DROP TABLE IF EXISTS events;
