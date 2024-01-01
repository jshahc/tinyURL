#!/usr/bin/env bash

# Function to print usage and exit
print_usage() {
    echo "Usage: $0 [CONFIG_FILE]"
    exit 1
}

# Function to read configuration and export variables
read_and_export_config() {
    local config_file=$1

    # Check if yq is installed
    if ! command -v yq &> /dev/null; then
        echo "Error: yq is not installed. Please install it before running this script."
        exit 1
    fi

    # Read configuration values from the specified config file
    PG_HOST=$(yq eval '.database.host' "$config_file")
    PG_PORT=$(yq eval '.database.port' "$config_file")
    PG_USER=$(yq eval '.database.username' "$config_file")
    PG_PASSWORD=$(yq eval '.database.password' "$config_file")
    PG_DB=$(yq eval '.database.dbname' "$config_file")

    # Server Configs
    STARTING_NUMBER=$(yq eval '.linkGenerator.startingNumber' "$config_file")
    GENERATOR_TYPE=$(yq eval '.linkGenerator.generatorType' "$config_file")

    # Export the variables
    export PG_HOST PG_PORT PG_USER PG_PASSWORD PG_DB STARTING_NUMBER GENERATOR_TYPE
}

# Function to create the PostgreSQL container
create_postgres_container() {
    echo "Starting PostgreSQL container with the following details:"
    echo "  Host: $PG_HOST"
    echo "  Port: $PG_PORT"
    echo "  User: $PG_USER"
    echo "  Password: $PG_PASSWORD" # This is not a good practice, but we are using it for simplicity
    echo "  Database: $PG_DB"
    echo ""

    # Startup a PostgreSQL container
    docker run --name my-postgres -e POSTGRES_PASSWORD="$PG_PASSWORD" -p "$PG_PORT":"$PG_PORT" -d postgres

    # Check the exit status of the Docker command
    local docker_status=$?
    if [ $docker_status -ne 0 ]; then
        echo "Error: Failed to start PostgreSQL container."
        exit $docker_status
    fi
}

# Function to create PostgreSQL tables
create_tables() {
    local psql_command="psql -h $PG_HOST -p $PG_PORT -U $PG_USER -d $PG_DB --password"

    # SQL command to create the short_links table with an index on the original_link column
    local sql_command="
    DROP TABLE IF EXISTS short_links;
    CREATE TABLE short_links (
        short_link VARCHAR (100) PRIMARY KEY,
        original_link VARCHAR (500) NOT NULL
    );
    CREATE INDEX original_link_index ON short_links (original_link);
    "

    echo "Creating the 'short_links' table in the database"

    # Execute SQL command using psql
    $psql_command -c "$sql_command"
    
    if [ "$GENERATOR_TYPE" = "sequential" ]; then
        local sql_command="
          DROP TABLE IF EXISTS counter_table;
          CREATE TABLE counter_table (
              counter_value BIGINT PRIMARY KEY
          );

          -- Insert  starting value into the table
          INSERT INTO counter_table (counter_value) VALUES ($STARTING_NUMBER);
        "

        echo "Creating the 'counter_table' table in the database"
        # Execute SQL command using psql
        $psql_command -c "$sql_command"
    fi
}

echo "Setting up PostgreSQL database using Docker and creating tables"

# Set the default configuration file
CONFIG_FILE="resources/config/local-run-config.yaml"

# Check if a configuration file parameter is provided
if [ "$#" -eq 1 ]; then
    CONFIG_FILE=$1
fi

echo "Using config file: $CONFIG_FILE"

# Read and export configuration variables
read_and_export_config "$CONFIG_FILE"

# Check if all required configuration values are present
if [ -z "$PG_HOST" ] || [ -z "$PG_PORT" ] || [ -z "$PG_USER" ] || [ -z "$PG_DB" ] || [ -z "$PG_PASSWORD" ]; then
    echo "Error: Configuration file is missing required values."
    print_usage
fi

# Create PostgreSQL container
create_postgres_container

# Create PostgreSQL tables
create_tables

echo "Setup completed successfully."
