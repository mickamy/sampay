CREATE TABLE user_profiles
(
    user_id    CHAR(26)                              NOT NULL PRIMARY KEY,
    name       VARCHAR(64)                           NOT NULL,
    bio        VARCHAR(256)                          NULL,
    image_id   CHAR(26)                              NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_profiles_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_user_profiles_image_id FOREIGN KEY (image_id) REFERENCES s3_objects (id)

);

CREATE INDEX idx_user_profiles_user_id ON user_profiles (user_id);
CREATE INDEX idx_user_profiles_name ON user_profiles (name);
