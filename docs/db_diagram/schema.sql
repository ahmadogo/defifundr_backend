-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2025-03-07T14:24:10.274Z

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "account_type" VARCHAR(50) NOT NULL,
  "personal_account_type" VARCHAR(50) NOT NULL,
  "first_name" VARCHAR(255) NOT NULL,
  "last_name" VARCHAR(255) NOT NULL,
  "nationality" VARCHAR(255) NOT NULL,
  "residencial_country" VARCHAR(255),
  "job_role" VARCHAR(255),
  "company_website" VARCHAR(255),
  "employment_type" VARCHAR(255),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "kyc" (
  "id" uuid PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "face_verification" bool NOT NULL,
  "identity_verification" bool NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "wallets" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" UUID NOT NULL,
  "wallet_address" VARCHAR(255) UNIQUE NOT NULL,
  "chain" VARCHAR(50) NOT NULL,
  "is_primary" BOOLEAN NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "pin_hash" VARCHAR(255) NOT NULL
);

CREATE TABLE "organizations" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" VARCHAR(255) NOT NULL,
  "employer_id" UUID NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "organization_members" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "organization_id" UUID NOT NULL,
  "employee_id" UUID NOT NULL,
  "role" VARCHAR(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "payrolls" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "employer_id" UUID NOT NULL,
  "organization_id" UUID,
  "payment_frequency" VARCHAR(50) NOT NULL,
  "salary_amount" "NUMERIC(18, 2)" NOT NULL,
  "currency" VARCHAR(10) NOT NULL,
  "contract_address" VARCHAR(255),
  "status" VARCHAR(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "payroll_employees" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "payroll_id" UUID NOT NULL,
  "employee_id" UUID NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "invoices" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "freelancer_id" UUID NOT NULL,
  "employer_id" UUID NOT NULL,
  "amount" "NUMERIC(18, 2)" NOT NULL,
  "currency" VARCHAR(10) NOT NULL,
  "status" VARCHAR(50) NOT NULL,
  "contract_address" VARCHAR(255),
  "created_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "transactions" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" UUID NOT NULL,
  "tx_hash" VARCHAR(255) UNIQUE NOT NULL,
  "amount" "NUMERIC(18, 2)" NOT NULL,
  "currency" VARCHAR(10) NOT NULL,
  "type" VARCHAR(50) NOT NULL,
  "status" VARCHAR(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "notifications" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" UUID NOT NULL,
  "message" TEXT NOT NULL,
  "type" VARCHAR(50) NOT NULL,
  "is_read" BOOLEAN NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE UNIQUE INDEX ON "users" ("email");

CREATE INDEX ON "wallets" ("user_id");

CREATE UNIQUE INDEX ON "wallets" ("wallet_address");

CREATE UNIQUE INDEX ON "organizations" ("name");

CREATE INDEX ON "organizations" ("employer_id");

CREATE INDEX ON "organization_members" ("organization_id");

CREATE INDEX ON "organization_members" ("employee_id");

CREATE UNIQUE INDEX ON "organization_members" ("organization_id", "employee_id");

CREATE INDEX ON "payrolls" ("employer_id");

CREATE INDEX ON "payrolls" ("organization_id");

CREATE INDEX ON "payroll_employees" ("payroll_id");

CREATE INDEX ON "payroll_employees" ("employee_id");

CREATE UNIQUE INDEX ON "payroll_employees" ("payroll_id", "employee_id");

CREATE INDEX ON "invoices" ("freelancer_id");

CREATE INDEX ON "invoices" ("employer_id");

CREATE INDEX ON "transactions" ("user_id");

CREATE UNIQUE INDEX ON "transactions" ("tx_hash");

CREATE INDEX ON "notifications" ("user_id");

CREATE INDEX ON "notifications" ("is_read");

COMMENT ON COLUMN "users"."account_type" IS 'business, personal';

COMMENT ON COLUMN "users"."personal_account_type" IS 'contractor, business';

COMMENT ON COLUMN "wallets"."chain" IS 'ethereum, solana';

COMMENT ON COLUMN "organization_members"."role" IS 'employee, manager, etc.';

COMMENT ON COLUMN "payrolls"."payment_frequency" IS 'weekly, bi-weekly, monthly';

COMMENT ON COLUMN "payrolls"."currency" IS 'USDC, SOL, ETH';

COMMENT ON COLUMN "payrolls"."status" IS 'pending, active, completed';

COMMENT ON COLUMN "invoices"."currency" IS 'USDC, SOL, ETH';

COMMENT ON COLUMN "invoices"."status" IS 'pending, approved, rejected, paid';

COMMENT ON COLUMN "transactions"."currency" IS 'USDC, SOL, ETH';

COMMENT ON COLUMN "transactions"."type" IS 'payroll, invoice';

COMMENT ON COLUMN "transactions"."status" IS 'pending, success, failed';

COMMENT ON COLUMN "notifications"."type" IS 'payroll, invoice, transaction';

ALTER TABLE "kyc" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wallets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "organizations" ADD FOREIGN KEY ("employer_id") REFERENCES "users" ("id");

ALTER TABLE "organization_members" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "organization_members" ADD FOREIGN KEY ("employee_id") REFERENCES "users" ("id");

ALTER TABLE "payrolls" ADD FOREIGN KEY ("employer_id") REFERENCES "users" ("id");

ALTER TABLE "payrolls" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "payroll_employees" ADD FOREIGN KEY ("payroll_id") REFERENCES "payrolls" ("id");

ALTER TABLE "payroll_employees" ADD FOREIGN KEY ("employee_id") REFERENCES "users" ("id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("freelancer_id") REFERENCES "users" ("id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("employer_id") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "notifications" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
