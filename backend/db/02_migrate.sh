#!/bin/bash

set -euo pipefail

echo "Starting database migration process..."

SQL_DIR="/docker-entrypoint-initdb.d/migrations"

execute_sql_file() {
  local file=$1
  echo "Executing migration file: $(basename "$file")..."
  if ! PGPASSWORD=$POSTGRES_PASSWORD psql -U "$POSTGRES_USER" -d "$DB_NAME" -f "$file"; then
    echo "Error: Failed to execute $(basename "$file")"
    return 1
  fi
  echo "Completed execution of $(basename "$file")"
}

if ! compgen -G "$SQL_DIR/*.up.sql" > /dev/null; then
  echo "Error: No .up.sql files found in $SQL_DIR"
  exit 1
fi

for sql_file in "$SQL_DIR"/*.up.sql; do
  if [ -f "$sql_file" ]; then
    if ! execute_sql_file "$sql_file"; then
      echo "Migration failed. Exiting..."
      exit 1
    fi
  else
    echo "Error: File $sql_file not found"
    exit 1
  fi
done

echo "Database migration completed successfully. All .up.sql migration files executed."
