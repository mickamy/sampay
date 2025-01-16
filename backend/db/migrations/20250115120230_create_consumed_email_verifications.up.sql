CREATE TABLE consumed_email_verifications
(
    email_verification_id CHAR(26)    NOT NULL PRIMARY KEY,
    consumed_at           TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_consumed_email_verifications_id FOREIGN KEY (email_verification_id) REFERENCES email_verifications (id)
);

CREATE INDEX idx_consumed_email_verifications_token ON consumed_email_verifications (consumed_at);
