#More information about how to run the software can be found on README.md

run:
	go run ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

up:
	docker-compose up -d

down:
	docker-compose down

wire:
	wire ./cmd/ordersystem

proto:
	protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

graphql:
	go run github.com/99designs/gqlgen generate