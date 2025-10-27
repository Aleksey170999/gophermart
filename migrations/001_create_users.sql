-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(63) NOT NULL,
    username VARCHAR(64) NOT NULL UNIQUE,
    current_balance INTEGER NOT NULL DEFAULT 0,
    withdrawn_balance INTEGER NOT NULL DEFAULT 0,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL
);
-- +goose Down
