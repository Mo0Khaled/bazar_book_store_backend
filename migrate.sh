#!/bin/bash

# Run the migrations
echo "Run migrations up"
goose postgres "${DB_URL}" up

# Start the app
go run ./main.go
