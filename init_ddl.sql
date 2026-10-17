CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    description VARCHAR(1024) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(256) NOT NULL,
    email VARCHAR(256) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Constraints are added conditionally so this script can be executed more than
-- once without failing on an already initialized database.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_users') THEN
        ALTER TABLE task
            ADD CONSTRAINT fk_users
            FOREIGN KEY (user_id) REFERENCES users(id);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'unique_email') THEN
        ALTER TABLE users
            ADD CONSTRAINT unique_email UNIQUE (email);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'unique_username') THEN
        ALTER TABLE users
            ADD CONSTRAINT unique_username UNIQUE (username);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_task_user_id ON task (user_id);
CREATE INDEX IF NOT EXISTS idx_task_created_at ON task (created_at);
