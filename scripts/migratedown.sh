#!/bin/bash

set -e

if [ -f .env ]; then
    echo "Loading .env..."
    source .env
else
    echo "No .env file found."
    exit 1
fi

if [ -z "$DB_URL" ]; then
    echo "DB_URL is not set. Exiting."
    exit 1
fi

echo "DB_URL=$DB_URL"

cd backend/sql/schema || exit 1

goose -v postgres "$DB_URL" down