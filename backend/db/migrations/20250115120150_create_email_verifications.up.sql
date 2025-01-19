CREATE TABLE email_verifications
(
    id          CHAR(26)                       NOT NULL PRIMARY KEY,
    intent_type email_verification_intent_type NOT NULL,
    email       VARCHAR(256)                   NOT NULL,
    created_at  TIMESTAMPTZ                    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_verifications_intent_type ON email_verifications (intent_type);
CREATE INDEX idx_email_verifications_email ON email_verifications (email);
