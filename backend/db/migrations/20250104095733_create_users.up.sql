CREATE TABLE users
(
    id         CHAR(26)     NOT NULL PRIMARY KEY,
    slug       VARCHAR(32)  NOT NULL,
    email      VARCHAR(256) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_users_slug UNIQUE (slug)
);

CREATE INDEX idx_users_slug ON users (slug);
CREATE INDEX idx_users_email ON users (email);
