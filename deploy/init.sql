CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    email VARCHAR,
    balance BIGINT NOT NULL CHECK (balance >= 0)
);

INSERT INTO users (name, email, balance) VALUES 
('Alice', 'alice@mail.ru', 1000), 
('Bob', 'bobmarley@gmail.com',2000);

