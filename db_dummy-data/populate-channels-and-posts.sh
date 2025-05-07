#!/usr/bin/env bash

set -euo pipefail

DB_PATH="../db/dev_forum_database.db"
SQL_FILE="populate_channels_and_posts.sql"
SQL_DELETE="delete_new.sql"

# echo "removing newlines from SQL file..."
# sqlite3 "$DB_PATH" <"$SQL_DELETE"
# echo "Done."
echo "Populating database..."
sqlite3 "$DB_PATH" <"$SQL_FILE"
echo "Done."
