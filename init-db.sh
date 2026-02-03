#!/bin/bash
set -e

# psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
#     CREATE DATABASE "$KC_DB_NAME";
# EOSQL

# This script runs against the API database (db-api)
# Use it to initialize extensions and tables
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    
    -- Insert your CREATE TABLE clients (...) logic here
    -- Or simply let your Go API handle migrations on startup
EOSQL