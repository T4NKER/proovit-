#!/bin/bash

set -e

bash ./scripts/databaseInit.sh
go mod download

go build -o proovit ./cmd

./proovit