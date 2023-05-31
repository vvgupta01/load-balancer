# go-load-balancer
Load balancer in Go that implements 
1. Random, round robin, weighted round robin, and least connections algorithms
2. Health check service
3. Client/server testing
## Local Run
```
go mod download
go run src/main.go
```
Run `go run src/main.go --help` for more info.
## Testing
```
go test ./test...
```
## Docker Run
```
./docker-run.sh
```