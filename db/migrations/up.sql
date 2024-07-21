DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_statuses') THEN
        CREATE TYPE message_statuses AS ENUM ('pending', 'processed');
    END IF;
END $$;

DROP TABLE messages;

CREATE TABLE IF NOT EXISTS messages(
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    status message_statuses DEFAULT 'pending'
);