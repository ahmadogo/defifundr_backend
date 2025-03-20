-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE invoices (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  freelancer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  employer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount NUMERIC(18, 2) NOT NULL,
  currency VARCHAR(10) NOT NULL,
  status VARCHAR(50) NOT NULL,
  contract_address VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01'
);

CREATE INDEX idx_invoices_freelancer_id ON invoices(freelancer_id);
CREATE INDEX idx_invoices_employer_id ON invoices(employer_id);

COMMENT ON COLUMN invoices.currency IS 'USDC, SOL, ETH';
COMMENT ON COLUMN invoices.status IS 'pending, approved, rejected, paid';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS invoices;