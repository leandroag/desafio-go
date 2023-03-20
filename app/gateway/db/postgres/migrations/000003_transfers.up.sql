BEGIN;

CREATE TABLE IF NOT EXISTS transfers (
    id SERIAL PRIMARY KEY,
    account_origin_id INTEGER NOT NULL REFERENCES accounts(id),
    account_destination_id INTEGER NOT NULL REFERENCES accounts(id),
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;