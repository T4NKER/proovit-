#!/bin/bash

# Exit on any error
set -e

# Fetch dependencies using go get
bash ./scripts/databaseInit.sh
go mod download

# Build the program
go build -o proovit-

# Run the built program
./proovit-