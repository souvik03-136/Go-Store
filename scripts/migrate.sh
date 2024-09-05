#!/bin/bash

# Exit script on error
set -e

# Database migration directory (replace with your actual path)
MIGRATION_DIR="./migrations"

# Database URL
DATABASE_URL="postgres://username:password@localhost:5432/dbname?sslmode=disable"

# Function to run migrations
run_migrations() {
    echo "Running database migrations..."

    # Apply migrations
    migrate -path $MIGRATION_DIR -database $DATABASE_URL up

    echo "Database migrations completed!"
}

# Function to rollback the last migration (if needed)
rollback_migration() {
    echo "Rolling back the last migration..."

    migrate -path $MIGRATION_DIR -database $DATABASE_URL down 1

    echo "Migration rollback completed!"
}

# Check arguments for rollback
if [ "$1" == "rollback" ]; then
    rollback_migration
else
    run_migrations
fi
