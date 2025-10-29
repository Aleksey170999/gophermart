-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    number VARCHAR(64) NOT NULL,
    status VARCHAR(63) NOT NULL,
    accural int NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    user_id int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose Down
