-- migrate:up
CREATE TABLE event_participants
(
    id         CHAR(26)    NOT NULL PRIMARY KEY,
    event_id   CHAR(26)    NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    tier       INT         NOT NULL DEFAULT 1,
    status     TEXT        NOT NULL DEFAULT 'unpaid',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_event_participants_event_id ON event_participants (event_id);

-- migrate:down
DROP TABLE IF EXISTS event_participants;
