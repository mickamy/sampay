CREATE TABLE user_links
(
    id            CHAR(26)                NOT NULL PRIMARY KEY,
    user_id       CHAR(26)                NOT NULL,
    provider_type user_link_provider_type NOT NULL,
    uri           VARCHAR(2048)           NOT NULL,
    created_at    TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_links_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_user_links_provider_type FOREIGN KEY (provider_type) REFERENCES user_link_providers (type)
);

CREATE INDEX idx_user_links_user_id ON user_links (user_id);
CREATE INDEX idx_user_links_provider_type ON user_links (provider_type);
