-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE payrolls (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  employer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
  payment_frequency VARCHAR(50) NOT NULL,
  salary_amount NUMERIC(18, 2) NOT NULL,
  currency VARCHAR(10) NOT NULL,
  contract_address VARCHAR(255),
  status VARCHAR(50) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01'
);

CREATE INDEX idx_payrolls_employer_id ON payrolls(employer_id);
CREATE INDEX idx_payrolls_organization_id ON payrolls(organization_id);

COMMENT ON COLUMN payrolls.payment_frequency IS 'weekly, bi-weekly, monthly';
COMMENT ON COLUMN payrolls.currency IS 'USDC, SOL, ETH';
COMMENT ON COLUMN payrolls.status IS 'pending, active, completed';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS payrolls;