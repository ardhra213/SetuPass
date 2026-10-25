#!/bin/bash
# Copyright 2025 The ThunderID Authors
# SPDX-License-Identifier: Apache-2.0

# Script common variables.
TYPE=""
SCHEMA_FILE_PATH=""
RECREATE="false"

# Database connection details.
HOST=""
PORT=""
NAME=""
DB_PATH=""
USERNAME=""
PASSWORD=""

print_help() {
  echo ""
  echo "Usage:"
  echo "  $0 [OPTIONS]"
  echo ""
  echo "Options:"
  printf "  %-12s %s\n" "-type" "Type of the database (e.g., postgres, sqlite) [required]"
  printf "  %-12s %s\n" "-schema" "Path to the database schema file [required]"
  printf "  %-12s %s\n" "-recreate" "Flag to recreate the database (no value needed)"
  printf "  %-12s %s\n" "-host" "Database host (required for postgres)"
  printf "  %-12s %s\n" "-port" "Database port (required for postgres)"
  printf "  %-12s %s\n" "-name" "Database name [required]"
  printf "  %-12s %s\n" "-path" "Path to the SQLite database file (required for sqlite)"
  printf "  %-12s %s\n" "-username" "Database username (required for postgres)"
  printf "  %-12s %s\n" "-password" "Database password (required for postgres)"
  printf "  %-12s %s\n" "-h, --help" "Show this help message and exit"
  echo ""
}

parse_args() {
  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      -type) TYPE="$2"; shift 2;;
      -schema) SCHEMA_FILE_PATH="$2"; shift 2;;
      -recreate) RECREATE="true"; shift;;
      -host) HOST="$2"; shift 2;;
      -port) PORT="$2"; shift 2;;
      -name) NAME="$2"; shift 2;;
      -path) DB_PATH="$2"; shift 2;;
      -username) USERNAME="$2"; shift 2;;
      -password) PASSWORD="$2"; shift 2;;
      -h|--help) print_help; exit 0;;
      *) echo "Unknown parameter passed: $1"; exit 1;;
    esac
  done
}

validate_inputs() {
  if [ -z "$SCHEMA_FILE_PATH" ]; then
    echo "Schema file path is required. Please provide it using -schema."
    echo "Use -h or --help for usage information."
    exit 1
  fi

  if [ ! -f "$SCHEMA_FILE_PATH" ]; then
    echo "Schema file not found: $SCHEMA_FILE_PATH"
    exit 1
  fi

  if [ -z "$TYPE" ]; then
    echo "Database type is required. Please provide it using -type."
    echo "Use -h or --help for usage information."
    exit 1
  fi

  if [ "$TYPE" = "postgres" ]; then
    if [ -z "$HOST" ] || [ -z "$PORT" ] || [ -z "$NAME" ] || [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
      echo "PostgreSQL connection details are required. Please provide -host, -port, -name, -username, and -password."
    echo "Use -h or --help for usage information."
      exit 1
    fi
  elif [ "$TYPE" = "sqlite" ]; then
    if [ -z "$DB_PATH" ]; then
      echo "SQLite database path is required. Please provide it using -path."
    echo "Use -h or --help for usage information."
      exit 1
    fi
  else
    echo "Unsupported database type: $TYPE"
    exit 1
  fi
}

init_postgres() {
  if [ "$RECREATE" = "true" ]; then
    echo "Dropping existing database..."
    PGPASSWORD="$PASSWORD" psql -h "$HOST" -p "$PORT" -U "$USERNAME" -d postgres -c "DROP DATABASE IF EXISTS $NAME;"
    [ $? -ne 0 ] && echo "Failed to drop database." && exit 1

    echo "Creating database..."
    PGPASSWORD="$PASSWORD" psql -h "$HOST" -p "$PORT" -U "$USERNAME" -d postgres -c "CREATE DATABASE $NAME;"
    [ $? -ne 0 ] && echo "Failed to create database." && exit 1
  else
    echo "Checking if database exists..."
    DB_EXISTS=$(PGPASSWORD="$PASSWORD" psql -h "$HOST" -p "$PORT" -U "$USERNAME" -d postgres -Atc "SELECT 1 FROM pg_database WHERE datname = '$NAME';" 2>/dev/null)

    if [ "$DB_EXISTS" != "1" ]; then
      echo "Database $NAME does not exist. Creating it..."
      PGPASSWORD="$PASSWORD" psql -h "$HOST" -p "$PORT" -U "$USERNAME" -d postgres -c "CREATE DATABASE $NAME;"
      [ $? -ne 0 ] && echo "Failed to create database." && exit 1
    else
      echo "Database $NAME already exists. Skipping creation."
    fi
  fi

  echo "Running schema on PostgreSQL..."
  PGPASSWORD="$PASSWORD" psql -h "$HOST" -p "$PORT" -U "$USERNAME" -d "$NAME" -f "$SCHEMA_FILE_PATH"
  [ $? -ne 0 ] && echo "Failed to run schema on PostgreSQL." && exit 1
}

init_sqlite() {
  DB_DIR=$(dirname "$DB_PATH")
  [ ! -d "$DB_DIR" ] && echo "Creating database directory: $DB_DIR" && mkdir -p "$DB_DIR"

  if [ "$RECREATE" = "true" ] && [ -f "$DB_PATH" ]; then
    echo "Removing existing SQLite database file..."
    rm -f "$DB_PATH" || { echo "Failed to remove existing database."; exit 1; }
  fi

  if [ ! -f "$DB_PATH" ]; then
    echo "Creating and initializing SQLite database..."
    sqlite3 "$DB_PATH" < "$SCHEMA_FILE_PATH"
    sqlite3 "$DB_PATH" "PRAGMA journal_mode=WAL;" # Enable WAL mode.
  else
    echo "SQLite database already exists. Skipping initialization."
  fi
}

init_database() {
  echo "Initializing the database..."
  case "$TYPE" in
    postgres) init_postgres;;
    sqlite) init_sqlite;;
    *) echo "Unsupported database type: $TYPE"; exit 1;;
  esac
}

main() {
  echo "Running database init script..."
  echo ""
  parse_args "$@"
  validate_inputs
  init_database
  echo ""
  echo "Database setup script completed."
}

main "$@"
