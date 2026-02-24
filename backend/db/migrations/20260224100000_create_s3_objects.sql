-- migrate:up
CREATE TABLE s3_objects
(
    id         CHAR(26)    NOT NULL PRIMARY KEY,
    bucket     TEXT        NOT NULL,
    key        TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (bucket, key)
);

-- migrate:down
DROP TABLE s3_objects;
