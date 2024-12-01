CREATE TABLE IF NOT EXISTS user_savings (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id VARCHAR PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    user_id VARCHAR NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    money NUMERIC(11, 2) NOT NULL DEFAULT 0,
    start_date TIMESTAMP
);