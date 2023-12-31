SHELL := /bin/bash
-include .env

.PHONY: prepare
prepare:
	cp .env.example .env
	cp .rabbit.env.example .rabbit.env
	cp .mysql.env.example .mysql.env

.PHONY: init
init:
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go generate ./internal/infra/grpc
	go generate ./internal/infra/gql
	go generate ./cmd/ordersystem/tools.go
	go mod tidy

.PHONY: build
build:
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/ordersystem

.PHONY: run
run:
	docker compose -f docker-compose.yaml up -d --build

.PHONY: migrate/up
migrate/up:
	migrate -path=sql/migrations \
		-database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(127.0.0.1:3306)/${DB_NAME}" \
		-verbose up

.PHONY: migrate/down
migrate/down:
	migrate -path=sql/migrations \
		-database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(127.0.0.1:3306)/${DB_NAME}" \
		-verbose down
