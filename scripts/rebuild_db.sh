#!/usr/bin/env bash

# Rebuild sCare MySQL data from modular seed SQL files.
# Usage:
#   ./scripts/rebuild_db.sh
#   MYSQL_CONTAINER=scare_mysql ./scripts/rebuild_db.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$ROOT_DIR"

if [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

MYSQL_CONTAINER="${MYSQL_CONTAINER:-scare_mysql}"
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:?DB_USER is required (from .env)}"
DB_PASSWORD="${DB_PASSWORD:?DB_PASSWORD is required (from .env)}"
DB_NAME="${DB_NAME:?DB_NAME is required (from .env)}"

SEED_FILES=(
  "backend/database/seeds/001_seed_permissions.sql"
  "backend/database/seeds/002_seed_users.sql"
  "backend/database/seeds/003_seed_stations.sql"
  "backend/database/seeds/004_seed_requests.sql"
  "backend/database/seeds/005_seed_notifications.sql"
  "backend/database/seeds/006_seed_news.sql"
  "backend/database/seeds/007_seed_menus.sql"
)

echo "Rebuilding database '$DB_NAME' in container '$MYSQL_CONTAINER'..."

for seed in "${SEED_FILES[@]}"; do
  echo "  -> applying $seed"
  docker exec -i "$MYSQL_CONTAINER" mysql \
    --default-character-set=utf8mb4 \
    -h"$DB_HOST" -P"$DB_PORT" \
    -u"$DB_USER" -p"$DB_PASSWORD" \
    "$DB_NAME" < "$seed"
done

echo "Database rebuild completed."
