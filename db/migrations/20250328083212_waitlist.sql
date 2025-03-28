-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE waitlist (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) NOT NULL UNIQUE,
  full_name VARCHAR(255),
  referral_code VARCHAR(20) NOT NULL UNIQUE,
  referral_source VARCHAR(100),
  status VARCHAR(20) NOT NULL DEFAULT 'waiting',
  signup_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  invited_date TIMESTAMPTZ,
  registered_date TIMESTAMPTZ,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_waitlist_status ON waitlist(status);
CREATE INDEX idx_waitlist_signup_date ON waitlist(signup_date);
CREATE UNIQUE INDEX idx_waitlist_referral_code ON waitlist(referral_code);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS waitlist;