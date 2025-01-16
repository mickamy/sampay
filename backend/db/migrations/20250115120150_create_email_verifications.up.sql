CREATE TABLE email_verifications
(
    id         CHAR(26)     NOT NULL PRIMARY KEY,
    email      VARCHAR(256) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_verification_email ON email_verifications (email);
