#!/bin/bash
# migrate_status.sh - Check the status of migrations

set -e

# Get directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MIGRATIONS_DIR="$DIR/../migrations"
DB_URL=${DATABASE_URL:-"postgres://postgres:postgres@localhost:5432/defifundr?sslmode=disable"}

# Check if goose is installed
if ! command -v goose &> /dev/null; then
    echo "Error: goose is not installed. Please install it with:"
    echo "go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

echo "Checking migration status"
echo "Migrations directory: $MIGRATIONS_DIR"
echo "Database URL: $DB_URL"

# Check migration status
goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" status

echo "Status check completed."