CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bank VARCHAR(3) NOT NULL CHECK(bank IN ('VCB', 'ACB', 'VIB')),
    balance FLOAT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT FK_AccountUser FOREIGN KEY (user_id)
    REFERENCES users(id) 
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    account_id SERIAL NOT NULL,
    amount FLOAT NOT NULL,
    transaction_type VARCHAR(10) NOT NULL CHECK(transaction_type IN ('withdraw', 'deposit')),
    created_at timestamp without time zone NOT NULL,
    CONSTRAINT FK_AccountTransaction FOREIGN KEY (account_id)
    REFERENCES accounts(id) 
);
-- CREATE UNIQUE INDEX idx_username ON customers(username);
