DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_statuses') THEN
        CREATE TYPE message_statuses AS ENUM ('pending', 'processed');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS messages(
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    status message_statuses DEFAULT 'pending'
);

CREATE INDEX IF NOT EXISTS idx_messages_id ON messages(id);
CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status);