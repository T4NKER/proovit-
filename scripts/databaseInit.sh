#!/bin/bash

# Get the directory of the script
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Define absolute paths based on the script directory
database_file="$script_dir/../src/database/database.db"
schema_file="$script_dir/../src/database/schema.sql"
seed_data_file="$script_dir/../src/database/seedData.sql"

# Check if the database file exists
if [ -e "$database_file" ]; then
    echo "Database file already exists. Doing nothing."
else
    # Create the database and apply schema
    sqlite3 "$database_file" < "$schema_file"
    sqlite3 "$database_file" < "$seed_data_file"
    echo "Database created with schema."
fi