#!/bin/bash

# Set your PostgreSQL connection details
PG_HOST="localhost"
PG_PORT="5432"
PG_USER="postgres"
PG_PASSWORD="mysecretpassword"
PG_DB="postgres"

# PSQL command to create the "original" table with an index on the "name" column
PSQL_COMMAND="psql -h $PG_HOST -p $PG_PORT -U $PG_USER -d $PG_DB --password $PG_PASSWORD"

# SQL command to create the "original" table with an index
SQL_COMMAND="
DROP TABLE IF EXISTS short_links;
CREATE TABLE short_links (
    short_link VARCHAR (100) PRIMARY KEY,
    original_link VARCHAR (500) NOT NULL
);
CREATE INDEX original_link_index ON short_links (original_link);
"

# Execute SQL command using psql
$PSQL_COMMAND -c "$SQL_COMMAND"
