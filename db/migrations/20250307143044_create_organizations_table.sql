-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE organizations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL UNIQUE,
  employer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01'
);

CREATE INDEX idx_organizations_employer_id ON organizations(employer_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS organizations;