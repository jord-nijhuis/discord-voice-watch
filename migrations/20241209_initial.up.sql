CREATE TABLE registrations(
    user_id TEXT NOT NULL,
    server_id TEXT NOT NULL,
    last_notified_at TIMESTAMP NULL,
    channel_id TEXT NULL,
    message_id TEXT NULL,
    PRIMARY KEY (user_id, server_id)
);

CREATE INDEX idx_registrations_server_id ON registrations (server_id)
