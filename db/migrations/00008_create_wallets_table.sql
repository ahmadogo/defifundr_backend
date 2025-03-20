-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE wallets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  wallet_address VARCHAR(255) NOT NULL UNIQUE,
  chain VARCHAR(50) NOT NULL,
  is_primary BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01',
  pin_hash VARCHAR(255) NOT NULL
);

CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_address ON wallets(wallet_address);

COMMENT ON COLUMN wallets.chain IS 'ethereum, solana';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS wallets;