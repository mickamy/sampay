CREATE TABLE user_profiles
(
    user_id    CHAR(26)                              NOT NULL PRIMARY KEY,
    name       VARCHAR(64)                           NOT NULL,
    bio        VARCHAR(256)                          NULL,
    image_id   CHAR(26)                              NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_profiles_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles (user_id);
CREATE INDEX idx_user_profiles_name ON user_profiles (name);
