CREATE TABLE notifications
(
    id         CHAR(26)          NOT NULL PRIMARY KEY,
    type       notification_type NOT NULL,
    user_id    CHAR(26)          NOT NULL,
    subject    VARCHAR(64)       NOT NULL,
    body       TEXT              NOT NULL,
    created_at TIMESTAMPTZ       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_notifications_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_notifications_type ON notifications (type);
CREATE INDEX idx_notifications_user_id ON notifications (user_id);
