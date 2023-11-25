# go-load-balancer
HTTP load balancer in Go that implements
1. Random, round robin, weighted round robin, and least connections algorithms
2. Concurrent health check service
3. REST API to view server config and server pool metrics
4. Unit tests, benchmarks, and client/server network simulations
## Local Run
```
// Download dependencies
go mod download
// Build & run executable
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
## REST API
```
// Returns server config/state
localhost:<MANAGER_PORT>/servers
// Returns server pool status/metrics
localhost:<MANAGER_PORT>/status
```