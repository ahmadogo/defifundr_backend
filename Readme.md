# DefiFundr - A decentralized crowdfunding platform for the Ethereum blockchain

[![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/demola234/deFICrowdFunding-Backend/test.yml)](https://github.com/DefiFundr-Labs/defifundr_backend/actions)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/demola234/deFICrowdFunding-Backend/main)
[![GitHub issues](https://img.shields.io/github/issues/demola234/deFICrowdFunding-Backend)](https://github.com/DefiFundr-Labs/defifundr_backend/issues?q=is%3Aissue%20state%3Aopen)
[![GitHub Repo stars](https://img.shields.io/github/stars/demola234/deFICrowdFunding-Backend)](https://github.com/DefiFundr-Labs/defifundr_backend/stargazers)

## What is DefiFundr?

DefiFundr is a revolutionary decentralized payroll and invoice management system that bridges the gap between traditional financial systems and blockchain technology. The platform provides a seamless, secure, and transparent solution for businesses to manage employee payments, handle freelancer invoices, and automate salary disbursements across both fiat and cryptocurrency channels.

## Installation

```bash
git clone
cd defifundr_backend
go mod download
```

## Usage

### Using Makefile

```bash
make server
```

### Using Go

```bash
go run main.go
```

### Using Air (Hot Reload)

```bash
air
```

## Testing

```bash
make test
```

### Unit Tests

```bash
go test ./...
```

### Coverage

```bash
go test -v -cover ./...
```

```markdown
# Database Migrations

The DefiFundr backend uses [goose](https://github.com/pressly/goose) for managing database migrations. Migrations are stored in the `migrations` directory and are written in SQL.

## Setting Up Migrations

1. Install goose:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. Make sure you have PostgreSQL running and accessible with the connection details specified in your environment variables or `.env` file.

## Migration Commands

We provide helper scripts to manage migrations:

* `./scripts/migrate.sh` - Apply all pending migrations
* `./scripts/migrate_create.sh <migration_name>` - Create a new migration file
* `./scripts/migrate_down.sh` - Roll back the last migration
* `./scripts/migrate_status.sh` - Check the status of all migrations
* `./scripts/migrate_reset.sh` - Roll back all migrations and apply them again (use with caution!)

### Creating a New Migration

To create a new migration:

```bash
./scripts/migrate_create.sh add_new_feature
```

This will create a new file in the `migrations` directory with a timestamp prefix, like `20230809123456_add_new_feature.sql`.

Edit this file to add your SQL commands. Each migration file should have two sections:

```sql
-- +goose Up
-- SQL in this section is executed when the migration is applied
CREATE TABLE example (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP TABLE example;
```

### Running Migrations

To apply all pending migrations:

```bash
./scripts/migrate.sh
```

### Database Connection

By default, the migration scripts will use the `DATABASE_URL` environment variable. If this is not set, they will fall back to `postgres://postgres:postgres@localhost:5432/defifundr?sslmode=disable`.

You can set the environment variable before running the script:

```bash
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable" ./scripts/migrate.sh
```

Or you can update the default value in the scripts.

## Migration Best Practices

1. **Always include a Down migration**: This ensures you can roll back if something goes wrong.
2. **Keep migrations small and focused**: It's better to have multiple small migrations than one large one.
3. **Test migrations before deploying**: Run migrations on a test database to ensure they work correctly.
4. **Version control migrations**: All migrations should be committed to the repository.
5. **Never modify an existing migration file**: Once a migration has been applied to any environment, create a new migration instead of modifying the existing one.

## Troubleshooting

If you encounter issues with migrations:

1. Check the migration status with `./scripts/migrate_status.sh`
2. Ensure your database connection details are correct
3. Look for syntax errors in your SQL statements
4. Check if you have the necessary permissions on the database

For more complex issues, you might need to manually fix the goose migration table (`goose_db_version`).
```

These scripts and documentation provide a comprehensive system for managing your database migrations using goose. Make sure to make the scripts executable after creating them (`chmod +x scripts/*.sh`).