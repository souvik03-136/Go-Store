#!/usr/bin/env bash
# scripts/migrate.sh
#
# Runs or rolls back database migrations using golang-migrate.
# Requires the `migrate` CLI: https://github.com/golang-migrate/migrate
#
# Usage:
#   ./scripts/migrate.sh          # apply all pending migrations
#   ./scripts/migrate.sh rollback # roll back the last applied migration
#   ./scripts/migrate.sh drop     # drop everything (DANGER: use only in dev)

set -euo pipefail

MIGRATION_DIR="./database/migrations"

# Build the DATABASE_URL from individual env vars when not already set.
if [[ -z "${DATABASE_URL:-}" ]]; then
  DB_DRIVER="${DB_DRIVER:-postgres}"
  DB_HOST="${DB_HOST:-localhost}"
  DB_PORT="${DB_PORT:-5432}"
  DB_DATABASE="${DB_DATABASE:-gostore}"
  DB_USERNAME="${DB_USERNAME:-postgres}"
  DB_PASSWORD="${DB_PASSWORD:-}"
  DB_SSL_MODE="${DB_SSL_MODE:-disable}"

  if [[ "$DB_DRIVER" == "mysql" ]]; then
    DATABASE_URL="mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}"
  else
    DATABASE_URL="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=${DB_SSL_MODE}"
  fi
fi

command -v migrate >/dev/null 2>&1 || {
  echo >&2 "ERROR: 'migrate' CLI not found. Install it:"
  echo >&2 "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
  exit 1
}

case "${1:-up}" in
  up)
    echo "Applying migrations…"
    migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" up
    echo "Migrations applied."
    ;;
  rollback|down)
    echo "Rolling back last migration…"
    migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" down 1
    echo "Rollback complete."
    ;;
  drop)
    echo "WARNING: dropping all tables in 5 seconds. Ctrl-C to abort."
    sleep 5
    migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" drop -f
    echo "All tables dropped."
    ;;
  *)
    echo "Usage: $0 [up|rollback|drop]"
    exit 1
    ;;
esac