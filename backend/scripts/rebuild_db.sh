#!/usr/bin/env bash

# Rebuild sCare MySQL data from modular seed SQL files.
# Usage:
#   ./scripts/rebuild_db.sh
#   MYSQL_CONTAINER=scare_mysql ./scripts/rebuild_db.sh

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
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
  "database/seeds/modules/00_reset_all.sql"
  "database/seeds/modules/10_roles_permissions.sql"
  "database/seeds/modules/20_menus.sql"
  "database/seeds/modules/30_stations_zones.sql"
  "database/seeds/modules/40_users_profiles.sql"
  "database/seeds/modules/50_requests_tasks.sql"
  "database/seeds/modules/60_content.sql"
  "database/seeds/modules/70_notifications.sql"
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
