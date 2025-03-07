#!/bin/bash
# migrate_reset.sh - Roll back all migrations and apply them again

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

echo "Resetting all migrations (down and up again)"
echo "Migrations directory: $MIGRATIONS_DIR"
echo "Database URL: $DB_URL"

# Confirm reset
read -p "This will roll back ALL migrations and apply them again. Continue? (y/N) " confirm
if [[ $confirm != [yY] && $confirm != [yY][eE][sS] ]]; then
    echo "Operation cancelled."
    exit 0
fi

# Roll back all migrations
echo "Rolling back all migrations..."
goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" down-to 0

# Apply all migrations
echo "Applying all migrations..."
goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" up

echo "Migration reset completed successfully!"