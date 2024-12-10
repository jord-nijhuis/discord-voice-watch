-- Create the users table
CREATE TABLE users
(
    id           TEXT PRIMARY KEY,
    channel_id   TEXT      NULL,
    snooze_until TIMESTAMP NULL
);


-- Insert the users from registrations into the new table
INSERT INTO users (id, channel_id)
SELECT user_id, channel_id
FROM registrations GROUP BY user_id;

CREATE TABLE servers
(
    id   TEXT PRIMARY KEY
);

INSERT INTO servers (id)
SELECT server_id
FROM registrations GROUP BY server_id;


-- Create a new registrations table
CREATE TABLE registrations_new
(
    user_id          TEXT      NOT NULL,
    server_id        TEXT      NOT NULL,
    last_notified_at TIMESTAMP NULL,
    message_id       TEXT      NULL,
    PRIMARY KEY (user_id, server_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (server_id) REFERENCES servers (id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- Copy the data from the old registrations table to the new one
INSERT INTO registrations_new (user_id, server_id, last_notified_at, message_id)
SELECT user_id, server_id, last_notified_at, message_id
FROM registrations;

-- Drop the old table
DROP TABLE registrations;

-- Rename the new table to the old one
ALTER TABLE registrations_new RENAME TO registrations;

-- Recreate the index
CREATE INDEX idx_registrations_server_id ON registrations (server_id);
