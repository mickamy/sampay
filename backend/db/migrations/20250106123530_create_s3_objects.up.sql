CREATE TABLE s3_objects
(
    id     CHAR(26)     NOT NULL PRIMARY KEY,
    bucket VARCHAR(64)  NOT NULL,
    key    VARCHAR(256) NOT NULL,
    CONSTRAINT uq_s3_objects_bucket_key UNIQUE (bucket, key)
);
