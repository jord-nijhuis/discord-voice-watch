-- Create a new registrations table --
CREATE TABLE registrations_new
(
    user_id          TEXT      NOT NULL,
    server_id        TEXT      NOT NULL,
    last_notified_at TIMESTAMP NULL,
    channel_id       TEXT      NULL,
    message_id       TEXT      NULL,
    PRIMARY KEY (user_id, server_id)
);

INSERT INTO registrations_new (user_id, server_id, last_notified_at, channel_id, message_id)
SELECT r.user_id, r.server_id, r.last_notified_at, u.channel_id, r.message_id
FROM registrations r
         LEFT JOIN user u ON r.user_id = u.id;

-- Drop the old table
DROP TABLE registrations;

-- Rename the new table to the old one
ALTER TABLE registrations_new RENAME TO registrations;

-- Recreate the index
CREATE INDEX idx_registrations_server_id ON registrations (server_id);

-- Drop the users table
DROP table users;

-- Drop the servers table
DROP TABLE servers;
