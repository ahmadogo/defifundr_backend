-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tx_hash VARCHAR(255) NOT NULL UNIQUE,
  amount NUMERIC(18, 2) NOT NULL,
  currency VARCHAR(10) NOT NULL,
  type VARCHAR(50) NOT NULL,
  status VARCHAR(50) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01'
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE UNIQUE INDEX idx_transactions_tx_hash ON transactions(tx_hash);

COMMENT ON COLUMN transactions.currency IS 'USDC, SOL, ETH';
COMMENT ON COLUMN transactions.type IS 'payroll, invoice';
COMMENT ON COLUMN transactions.status IS 'pending, success, failed';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS transactions;