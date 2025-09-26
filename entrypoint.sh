#!/bin/sh
set -e

DB_FILE=/var/lib/db-codex/forum_database.db

# Initialize DB if missing
if [ ! -f "$DB_FILE" ]; then
  echo "Initializing database..."
  /app/codex migrate /migrations/001_schema.sql
  /app/codex seed
fi

exec "$@"
