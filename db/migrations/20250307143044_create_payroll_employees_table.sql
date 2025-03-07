-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE payroll_employees (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  payroll_id UUID NOT NULL REFERENCES payrolls(id) ON DELETE CASCADE,
  employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01'
);

CREATE INDEX idx_payroll_employees_payroll_id ON payroll_employees(payroll_id);
CREATE INDEX idx_payroll_employees_employee_id ON payroll_employees(employee_id);
CREATE UNIQUE INDEX idx_payroll_employees_unique ON payroll_employees(payroll_id, employee_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS payroll_employees;