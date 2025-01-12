CREATE TABLE user_link_display_attributes
(
    user_link_id  CHAR(26)     NOT NULL PRIMARY KEY,
    name          VARCHAR(256) NOT NULL,
    display_order SMALLINT     NOT NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_link_display_attributes_user_link_id FOREIGN KEY (user_link_id) REFERENCES user_links (id),
    CONSTRAINT uq_user_link_display_attributes_display_order UNIQUE (display_order) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_user_link_display_attributes_user_link_id ON user_link_display_attributes (user_link_id);
