#!/bin/bash

# Fix Go module cache & sumdb permission issues
# and lock gorilla/sessions to v1.3.0-compatible setup

set -e  # Exit on any error

echo "→ Setting environment variables..."
export GOSUMDB=off
export GOMODCACHE=$HOME/gomodcache

echo "→ Creating clean module cache..."
mkdir -p "$GOMODCACHE"

echo "→ Cleaning old module cache..."
go clean -modcache

echo "→ Tidying go.mod and downloading dependencies..."
go mod tidy

echo "→ Ensuring Middlewares package is resolved..."
go get parkplace/Middlewares

echo ""
echo "✅ Done! You can now run:"
echo "   go run main.go"
echo ""
echo "Current sessions version:"
go list -m github.com/gorilla/sessions
