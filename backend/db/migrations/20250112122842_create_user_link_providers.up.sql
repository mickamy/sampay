CREATE TABLE user_link_providers
(
    type          user_link_provider_type NOT NULL PRIMARY KEY,
    display_order SMALLINT                NOT NULL,
    created_at    TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_link_providers_display_order UNIQUE (display_order) DEFERRABLE INITIALLY DEFERRED
);
