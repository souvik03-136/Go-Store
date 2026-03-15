#!/usr/bin/env bash
# scripts/deploy.sh
#
# Deploys the application. With Docker Compose it rebuilds and restarts
# containers; without it it does a direct binary swap.

set -euo pipefail

echo "==> Pulling latest code…"
git pull origin main

echo "==> Running database migrations…"
./scripts/migrate.sh up

if [[ -f "docker-compose.yml" ]]; then
  echo "==> Docker Compose detected — rebuilding and restarting…"
  docker compose down --remove-orphans
  docker compose build --no-cache
  docker compose up -d
  echo "==> Deployment complete. Containers running:"
  docker compose ps
else
  echo "==> Building binary…"
  go build -o bin/api ./cmd/api

  if pgrep -f "bin/api" > /dev/null; then
    echo "==> Stopping existing process…"
    pkill -f "bin/api"
    sleep 2
  fi

  echo "==> Starting new binary…"
  nohup ./bin/api >> logs/app.log 2>&1 &
  echo "==> Deployment complete. PID: $!"
fi