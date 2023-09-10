gen:
	protoc --proto_path=./proto  --go_out=plugins=grpc:pb  proto/*.proto

clean:
	rm pb/*.go

server:
	go run cmd/server/main.go -port 8080
client:
	go run cmd/client/main.go -address 0.0.0.0:8080

test:
	go env -w CGO_ENABLED=1
	go test -cover -race ./...
cert:
	cd cert; ./gen.sh; cd ..

.PHONY: gen clean server client test cert