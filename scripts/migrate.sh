#!/bin/bash
# migrate.sh - Run migrations up to the latest version

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

echo "Running migrations from: $MIGRATIONS_DIR"
echo "Database URL: $DB_URL"

# Run migrations
goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" up

echo "Migrations completed successfully!"