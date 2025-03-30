# DefiFundr Database Guide

This document provides comprehensive information about the DefiFundr database architecture, schema, migration workflows, and query patterns.

## Table of Contents

- [DefiFundr Database Guide](#defifundr-database-guide)
  - [Table of Contents](#table-of-contents)
  - [Database Overview](#database-overview)
  - [Schema](#schema)
    - [Generating Schema Documentation](#generating-schema-documentation)
  - [Migration Workflow](#migration-workflow)
    - [Creating a New Migration](#creating-a-new-migration)
    - [Applying Migrations](#applying-migrations)
    - [Reverting Migrations](#reverting-migrations)
    - [Checking Migration Status](#checking-migration-status)
  - [SQL Generation with SQLC](#sql-generation-with-sqlc)
    - [Workflow](#workflow)
    - [Example SQL Query](#example-sql-query)
  - [Query Patterns](#query-patterns)
    - [Repository Pattern](#repository-pattern)
    - [Transaction Management](#transaction-management)
  - [Database Seeding](#database-seeding)
  - [Best Practices](#best-practices)

## Database Overview

DefiFundr uses PostgreSQL as its primary database. The database stores user information, authentication data, transactions, and other application state. We use the following tools to manage our database:

- **PostgreSQL** - Relational database
- **Goose** - Database migration management
- **SQLC** - SQL compiler for Go
- **DBML** - Database markup language for documentation

## Schema

The database schema includes these primary tables:

- `users` - User accounts and profile information
- `sessions` - User authentication sessions
- `kyc` - Know Your Customer verification data
- `user_device` - User device information
- `otp` - One-time password verification records
- `transactions` - Financial transaction records

The complete schema is defined in migration files located in `db/migrations/`. You can also view the schema in DBML format at `docs/db_diagram/db.dbml`.

### Generating Schema Documentation

To generate updated schema documentation:

```bash
# Generate SQL schema from DBML
make db_schema

# Generate DB documentation
make db_docs
```

## Migration Workflow

We use [Goose](https://github.com/pressly/goose) for database migrations. Migration files are stored in `db/migrations/` and follow a sequential numbering pattern.

### Creating a New Migration

To create a new migration:

```bash
make migrate-create
# Enter a descriptive name when prompted
```

This creates a new SQL migration file with up and down migrations.

### Applying Migrations

To apply pending migrations:

```bash
# Apply all pending migrations
make migrate-up

# Apply only the next pending migration
make migrate-up-one
```

### Reverting Migrations

To revert migrations:

```bash
# Revert the most recent migration
make migrate-down

# Revert the most recent migration only
make migrate-down-one

# Revert all migrations
make migrate-reset
```

### Checking Migration Status

To see the current migration status:

```bash
make migrate-status
```

## SQL Generation with SQLC

We use [SQLC](https://github.com/kyleconroy/sqlc) to generate type-safe Go code from SQL queries. SQL queries are defined in `db/query/` directory.

### Workflow

1. Define your SQL queries in files under `db/query/`
2. Run the SQL code generator:

```bash
make sqlc
```

3. Use the generated Go code in your repositories

SQLC generates:
- Strong types for rows and parameters
- Idiomatic Go functions for each query
- Interface definitions for mocking

### Example SQL Query

Here's an example SQL query from `db/query/users.sql`:

```sql
-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  name, email, hashed_password, created_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;
```

## Query Patterns

### Repository Pattern

We implement the repository pattern using interfaces defined in `internal/core/ports/repository.go` and implemented in `internal/adapters/repositories/`.

### Transaction Management

For operations that require multiple queries in a transaction, use the `Store` interface that provides transaction support.

Example:
```go
err := store.ExecTx(ctx, func(q *Queries) error {
    // Execute multiple queries within a transaction
    return nil
})
```

## Database Seeding

For development and testing, you can seed the database with sample data:

```bash
make seed
```

The seeding logic is defined in `cmd/seed/main.go`.

## Best Practices

1. **Migration Safety**:
   - Always include both "up" and "down" migrations
   - Test migrations in development before applying to production
   - Avoid modifying existing migrations after they've been applied

2. **Query Organization**:
   - Group related queries in the same .sql file
   - Use clear, descriptive query names
   - Document complex queries with comments

3. **Performance**:
   - Add appropriate indexes for frequently queried columns
   - Use EXPLAIN ANALYZE to check query performance
   - Implement pagination for large result sets using the utilities in `pkg/pagination`

4. **Testing**:
   - Write tests for repository implementations
   - Use the mock repository for service-level testing

5. **Schema Changes**:
   - Document significant schema changes in PR descriptions
   - Consider data migration needs when changing schemas
   - Use database constraints to enforce data integrity