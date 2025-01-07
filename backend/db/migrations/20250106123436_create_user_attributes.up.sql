CREATE TABLE user_attributes
(
    user_id             CHAR(26)            NOT NULL PRIMARY KEY,
    usage_category_type usage_category_type NOT NULL,
    created_at          TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_attributes_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_user_attributes_user_id ON user_attributes (user_id);
CREATE INDEX idx_user_attributes_usage_category_type ON user_attributes (usage_category_type);
