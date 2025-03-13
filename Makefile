test:
	@go test ./... -v

bench:
	@go test -bench=. ./...

generate:
	@protoc --go_out=. --go-grpc_out=. event.proto
