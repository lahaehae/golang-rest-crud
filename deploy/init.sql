CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    email VARCHAR,
    balance BIGINT NOT NULL CHECK (balance >= 0)
);

INSERT INTO users (name, email, balance) VALUES 
('Alice', 'alice@mail.ru', 1000), 
('Bob', 'bobmarley@gmail.com',2000);

CREATE TABLE sub_users (
    id UUID PRIMARY KEY,
    owner_id VARCHAR(36) NOT NULL,  -- ID пользователя из Kratos
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(owner_id, email)
);

CREATE INDEX idx_sub_users_owner ON sub_users(owner_id);
