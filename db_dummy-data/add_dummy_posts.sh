#!/bin/sh

DB_PATH="../db/forum_database.db"

run_sqlite_command() {
  sql_file=$1
  task_name=$2
  echo "Starting: $task_name"
  if sqlite3 "$DB_PATH" <"$sql_file"; then
    echo "$task_name completed successfully."
  else
    echo "Error: Failed to execute $task_name." >&2
    exit 1
  fi
}

run_sqlite_command "insert_posts.sql" "Creating posts"
run_sqlite_command "insert_postchannels.sql" "Creating channel references"
run_sqlite_command "insert_avatars.sql" "Inserting avatars"

echo "All tasks completed successfully."
