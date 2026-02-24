-- migrate:up
CREATE TABLE user_payment_methods
(
    id                  CHAR(26)    NOT NULL PRIMARY KEY,
    user_id             CHAR(26)    NOT NULL REFERENCES end_users (user_id),
    type                TEXT        NOT NULL,
    url                 TEXT        NOT NULL,
    qr_code_s3_object_id CHAR(26)  REFERENCES s3_objects (id),
    display_order       INT         NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, type)
);
CREATE INDEX idx_user_payment_methods_user_id ON user_payment_methods (user_id);

-- migrate:down
DROP TABLE IF EXISTS user_payment_methods;
