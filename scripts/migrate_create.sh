#!/bin/bash
# migrate_create.sh - Create a new migration file

set -e

# Get directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MIGRATIONS_DIR="$DIR/../db/migrations"

# Check if goose is installed
if ! command -v goose &> /dev/null; then
    echo "Error: goose is not installed. Please install it with:"
    echo "go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

# Check if a name was provided
if [ -z "$1" ]; then
    echo "Error: No migration name provided"
    echo "Usage: $0 <migration_name>"
    exit 1
fi

MIGRATION_NAME="$1"

echo "Creating migration: $MIGRATION_NAME"
goose -dir "$MIGRATIONS_DIR" create "$MIGRATION_NAME" sql

echo "Migration file created successfully in $MIGRATIONS_DIR"