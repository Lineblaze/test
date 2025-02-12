CREATE TABLE coin_transactions (
    id SERIAL PRIMARY KEY,
    from_user TEXT REFERENCES users(username) ON DELETE SET NULL,
    to_user TEXT REFERENCES users(username) ON DELETE SET NULL,
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);