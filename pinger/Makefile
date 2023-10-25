.PHONY: generate
generate:

		mkdir -p pkg/reporter/service
		protoc --go_out=pkg/reporter/service --go_opt=paths=import \
				--go-grpc_out=pkg/reporter/service --go-grpc_opt=paths=import \
				api/reporter/service.proto
		mv pkg/reporter/service/reporter/pkg/reporter/service/* pkg/reporter/service
		rm -rf pkg/reporter/service/reporter

.PHONY: build
build:
		go build -o bin/client cmd/client.go

.PHONY: run
run:
		cd bin && ./client

.PHONY: ping
ping: build run