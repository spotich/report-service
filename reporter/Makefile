PHONY: generate
generate:
		mkdir -p pkg/timer/service
		protoc --go_out=pkg/timer/service --go_opt=paths=import \
					--go-grpc_out=pkg/timer/service --go-grpc_opt=paths=import \
					api/timer/service.proto
		mv pkg/timer/service/reporter/pkg/timer/service/* pkg/timer/service
		rm -rf pkg/timer/service/reporter

		mkdir -p pkg/reporter/service
		protoc --go_out=pkg/reporter/service --go_opt=paths=import \
				--go-grpc_out=pkg/reporter/service --go-grpc_opt=paths=import \
				api/reporter/service.proto
		mv pkg/reporter/service/reporter/pkg/reporter/service/* pkg/reporter/service
		rm -rf pkg/reporter/service/reporter

.PHONY: build
build:
		go build -o bin/server cmd/server.go

.PHONY: run
run:
		cd bin && ./server

.PHONY: start
start: build run