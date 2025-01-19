CREATE TABLE verified_email_verifications
(
    email_verification_id CHAR(26)    NOT NULL PRIMARY KEY,
    token                 CHAR(64)    NOT NULL,
    verified_at           TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_verified_email_verifications_id FOREIGN KEY (email_verification_id) REFERENCES email_verifications (id)
);

CREATE INDEX idx_verified_email_verifications_token ON verified_email_verifications (token);
CREATE INDEX idx_verified_email_verifications_verified_at ON verified_email_verifications (verified_at);
