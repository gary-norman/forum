#!/usr/bin/env bash

set -euo pipefail

DB_PATH="/var/lib/db-codex/dev_forum_database.db"
SQL_FILE="insert_posts_and_channels.sql"

echo "Populating database..."
sqlite3 "$DB_PATH" <"$SQL_FILE"
echo "Done."
