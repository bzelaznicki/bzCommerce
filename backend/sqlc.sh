#!/bin/bash
export $(grep -v '^#' .env | xargs)
sqlc generate
