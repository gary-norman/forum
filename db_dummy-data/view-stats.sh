#!/usr/bin/env bash

set -euo pipefail

DB_PATH="/var/lib/db-codex/dev_forum_database.db"
SQL_STATS="view_stats.sql"

echo "Querying database..."
sqlite3 "$DB_PATH" <"$SQL_STATS"
echo "Done."
