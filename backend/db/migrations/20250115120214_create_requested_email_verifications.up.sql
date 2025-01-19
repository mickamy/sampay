CREATE TABLE requested_email_verifications
(
    email_verification_id CHAR(26)    NOT NULL PRIMARY KEY,
    pin_code              CHAR(6)     NOT NULL,
    requested_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at            TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_requested_email_verifications_id FOREIGN KEY (email_verification_id) REFERENCES email_verifications (id)
);

CREATE INDEX idx_requested_email_verifications_pin_code ON requested_email_verifications (pin_code);
CREATE INDEX idx_requested_email_verifications_expires_at ON requested_email_verifications (expires_at);
