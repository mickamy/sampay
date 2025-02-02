CREATE TABLE notifications
(
    id         CHAR(26)    NOT NULL PRIMARY KEY,
    user_id    CHAR(26)    NOT NULL,
    subject    VARCHAR(64) NOT NULL,
    body       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_notifications_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_notifications_user_id ON notifications (user_id);
