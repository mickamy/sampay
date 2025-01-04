CREATE TABLE authentications
(
    id          CHAR(26)                 NOT NULL PRIMARY KEY,
    user_id     CHAR(26)                 NOT NULL,
    type        authentication_type_enum NOT NULL,
    identifier  VARCHAR(256)             NOT NULL,
    secret      TEXT                     NULL,
    mfa_enabled BOOLEAN                  NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_authentications_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT uq_authentications_user_id_identifier UNIQUE (user_id, identifier)
);

CREATE INDEX idx_authentications_user_id ON authentications (user_id);
CREATE INDEX idx_authentications_type_identifier ON authentications (type, identifier);
