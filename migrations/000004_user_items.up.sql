CREATE TABLE user_items(
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL UNIQUE,
    quantity INT DEFAULT 0,
    UNIQUE (user_id, type)
);

