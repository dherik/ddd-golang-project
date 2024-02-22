
INSERT INTO users (username, email, password) VALUES
    ('admin', 'admin@example.com', '$2a$08$451vrF7qAb0lpcIUHeHbW.mFQnpCTgXoFUUlnskC9X6FYLywsN//G'),
    ('john_doe', 'john.doe@example.com', '$2a$08$451vrF7qAb0lpcIUHeHbW.mFQnpCTgXoFUUlnskC9X6FYLywsN//G'),
    ('jane_smith', 'jane.smith@example.com', '$2a$08$451vrF7qAb0lpcIUHeHbW.mFQnpCTgXoFUUlnskC9X6FYLywsN//G');

INSERT INTO Task (user_id, description, created_at) VALUES
    (1, 'Complete project proposal', '2024-02-15 10:59:01.054'),
    (2, 'Review code for bug fixes', '2024-02-16 21:23:09.066');