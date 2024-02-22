CREATE TABLE IF NOT EXISTS Task (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE Task
    ADD CONSTRAINT fk_users
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE users
    ADD CONSTRAINT unique_email UNIQUE (email);
