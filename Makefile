run:
	go run ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

up:
	docker-compose up -d

down:
	docker-compose down

wire:
	wire ./cmd/ordersystem