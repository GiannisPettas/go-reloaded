@echo off
echo Clearing all Go caches...
go clean -cache -testcache -modcache -fuzzcache

echo Running all tests without cache...
cd internal\testutils
go test -count=1 -v -run TestAllProject