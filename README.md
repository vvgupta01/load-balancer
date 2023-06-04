# go-load-balancer
Load balancer in Go that implements
1. Random, round robin, weighted round robin, and least connections strategies
2. Concurrent health check service
3. Unit tests, benchmarks, and client/server network testing
## Local Run
```
go mod download
go run src/main.go
```
Run `go run src/main.go --help` for more info.
## Testing
```
// Unit tests
go test ./test...
// Include benchmarks
go test ./test... -bench .
```
## Docker Run
```
./docker-run.sh
```