#!/usr/bin/env bash

set -euo pipefail

echo

echo "*** CLEANING UP CABIN-CHAT ***"

echo

echo "Formatting Go project..."
go fmt ./...

echo

echo "*** FORMATTING COMPLETED ***"

echo

echo "Building Go project..."
go build ./...

echo

echo "*** BUILD COMPLETED ***"

echo 

echo "Running Go tests..."
GOTMPDIR="$PWD/.tmp/go-build" go test -p 1 -v ./...

echo

echo "*** TESTS COMPLETED ***"

echo

echo "All done!"