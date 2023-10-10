PHONY: generate
generate:
		mkdir -p pkg/service
		protoc --go_out=pkg/service --go_opt=paths=import \
				--go-grpc_out=pkg/service --go-grpc_opt=paths=import \
				api/service.proto
		mv pkg/service/pinger/pkg/service/* pkg/service
		rm -rf pkg/service/pinger