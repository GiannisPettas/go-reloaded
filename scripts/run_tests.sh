#!/bin/bash

# Run all tests for go-reloaded project
echo "Running all tests..."
cd "$(dirname "$0")/.."

# Run tests and capture output (force fresh run)
if go test -count=1 ./...; then
    echo "✅ All tests passed!"
else
    echo "❌ Tests failed! See output above for details."
    exit 1
fi