#!/bin/bash
echo "Clearing all Go caches..."
go clean -cache -testcache -modcache -fuzzcache

echo "Running Golden Tests without cache..."
cd internal/testutils
go test -count=1 -v -run TestGoldenCases