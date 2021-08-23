clean:
	rm pb/*.go \
	rm -rf tmp/*

gen:
	protoc --proto_path=proto proto/*.proto --go_out=./pb --go-grpc_out=./pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

run:
	go run main.go

test:
	unset https_proxy http_proxy all_proxy; \
	go test -cover ./...;

server:
	go run cmd/server/main.go -port 9000

client:
	go run cmd/client/main.go -address 0.0.0.0:9000

.PHONY: clean gen run test client