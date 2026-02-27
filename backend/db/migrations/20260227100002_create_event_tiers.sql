-- migrate:up
CREATE TABLE event_tiers
(
    id         CHAR(26)    NOT NULL PRIMARY KEY,
    event_id   CHAR(26)    NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    tier       INT         NOT NULL,
    count      INT         NOT NULL,
    amount     INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (event_id, tier)
);

CREATE INDEX idx_event_tiers_event_id ON event_tiers (event_id);

-- migrate:down
DROP TABLE IF EXISTS event_tiers;
