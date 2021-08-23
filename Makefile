clean:
	rm pb/*.go \
	rm -rf tmp/*

gen:
	protoc --proto_path=proto proto/*.proto --go_out=./pb --go-grpc_out=./pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

run:
	go run main.go

test:
	go test -v -cover ./...


.PHONY: clean gen run test