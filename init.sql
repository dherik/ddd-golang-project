-- init.sql

CREATE TABLE IF NOT EXISTS Task (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "User" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE Task
    ADD CONSTRAINT fk_user
    FOREIGN KEY (user_id) REFERENCES "User"(id);

ALTER TABLE "User"
    ADD CONSTRAINT unique_email UNIQUE (email);

INSERT INTO "User" (username, email) VALUES
    ('john_doe', 'john.doe@example.com'),
    ('jane_smith', 'jane.smith@example.com');

INSERT INTO Task (user_id, description) VALUES
    ("1", 'Complete project proposal'),
    ("2", 'Review code for bug fixes');