package grpc

//go:generate protoc --go_out=. --go-grpc_out=. ./protofiles/order.proto

import (
	_ "github.com/kenesparta/fullcycle-clean-architecture/internal/infra/grpc/pb"
	_ "github.com/kenesparta/fullcycle-clean-architecture/internal/infra/grpc/service"
)
