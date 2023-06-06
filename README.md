# go-load-balancer
Load balancer in Go that implements
1. Random, round robin, weighted round robin, and least connections strategies
2. Concurrent health check service
3. REST API to view server config and server pool metrics
4. Unit tests, benchmarks, and client/server network testing
## Local Run
```
// Download dependencies
go mod download

// Build executable
go build -o loadbalancer ./src/main

./loadbalancer
```
Run `./loadbalancer --help` for more info.
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