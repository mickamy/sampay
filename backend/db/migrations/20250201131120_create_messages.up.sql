CREATE TABLE messages
(
    id          CHAR(26)     NOT NULL PRIMARY KEY,
    sender_name VARCHAR(256) NOT NULL,
    receiver_id CHAR(26)     NOT NULL,
    content     TEXT         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_messages_receiver_id FOREIGN KEY (receiver_id) REFERENCES users (id)
);

CREATE INDEX idx_messages_receiver_id ON messages (receiver_id);
CREATE INDEX idx_messages_created_at_desc ON messages (created_at DESC);
