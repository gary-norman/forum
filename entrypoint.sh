#!/bin/sh
set -e

DB_FILE=/var/lib/db-codex/forum_database.db

# Initialize DB if missing
if [ ! -f "$DB_FILE" ]; then
  echo "Initializing database..."
  /app/bin/codex migrate
  echo "Seeding database..."
  /app/bin/codex seed
fi

exec "$@"
