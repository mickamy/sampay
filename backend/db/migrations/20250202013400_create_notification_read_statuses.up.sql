CREATE TABLE notification_read_statuses
(
    notification_id CHAR(26)    NOT NULL,
    user_id         CHAR(26)    NOT NULL,
    read_at         TIMESTAMPTZ NULL,
    PRIMARY KEY (notification_id, user_id),
    CONSTRAINT fk_notification_read_status_notification_id FOREIGN KEY (notification_id) REFERENCES notifications (id),
    CONSTRAINT fk_notification_read_status_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_notification_read_statuses_notification_id ON notification_read_statuses (notification_id);
CREATE INDEX idx_notification_read_statuses_user_id ON notification_read_statuses (user_id);
CREATE INDEX idx_notification_read_statuses_read_at ON notification_read_statuses (read_at);
