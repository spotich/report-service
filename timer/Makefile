.PHONY: generate
generate:
		mkdir -p pkg/service
		protoc --go_out=pkg/service --go_opt=paths=import \
				--go-grpc_out=pkg/service --go-grpc_opt=paths=import \
				api/service.proto
		mv pkg/service/timer/pkg/timer/service/* pkg/service
		rm -rf pkg/service/timer

.PHONY: build
build:
		go build -o bin/server cmd/server.go

.PHONY: run
run:
		cd bin && ./server

.PHONY: start
start: build run
