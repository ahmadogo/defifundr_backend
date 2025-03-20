#!/bin/bash
# migrate_down.sh - Roll back the last migration

set -e

# Get directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MIGRATIONS_DIR="$DIR/../db/migrations"
DB_URL=${DATABASE_URL:-"postgres://root:secret@localhost:5433/defi?sslmode=disable"}

# Check if goose is installed
if ! command -v goose &> /dev/null; then
    echo "Error: goose is not installed. Please install it with:"
    echo "go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

echo "Rolling back the last migration"
echo "Migrations directory: $MIGRATIONS_DIR"
echo "Database URL: $DB_URL"

# Roll back one migration
goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" down

echo "Migration rollback completed successfully!"