SHELL := /bin/bash

init:
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go generate ./internal/infra/grpc
	go generate ./internal/infra/gql
	go generate ./cmd/ordersystem/wire.go
	go mod tidy

prod/build:
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/ordersystem

prod: init prod/build

prod/dkr:
	docker-compose -f docker-compose.prod.yml up -d --build

dev/dkr:
	docker-compose -f docker-compose.dev.yml up -d --build
