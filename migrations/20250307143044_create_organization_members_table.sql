-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE organization_members (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role VARCHAR(50) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01'
);

CREATE INDEX idx_org_members_org_id ON organization_members(organization_id);
CREATE INDEX idx_org_members_employee_id ON organization_members(employee_id);
CREATE UNIQUE INDEX idx_org_members_unique ON organization_members(organization_id, employee_id);

COMMENT ON COLUMN organization_members.role IS 'employee, manager, etc.';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS organization_members;