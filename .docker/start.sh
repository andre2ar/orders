#!/bin/bash

if [ ! -f "./cmd/ordersystem/app.env" ]; then
  echo "Creating app.env file..."
  cp ./cmd/ordersystem/app.env.example ./cmd/ordersystem/app.env
fi

go mod tidy

echo "Waiting to connect to RabbitMQ..."
sleep 3

echo "Starting server..."
go run ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go