package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kenesparta/fullcycle-clean-architecture/config"
	"github.com/kenesparta/fullcycle-clean-architecture/internal/event/handler"
	"github.com/kenesparta/fullcycle-clean-architecture/internal/infra/gql"
	"github.com/kenesparta/fullcycle-clean-architecture/internal/infra/gql/graph"
	"github.com/kenesparta/fullcycle-clean-architecture/internal/infra/grpc/pb"
	"github.com/kenesparta/fullcycle-clean-architecture/internal/infra/grpc/service"
	"github.com/kenesparta/fullcycle-clean-architecture/internal/infra/web/webserver"
	"github.com/kenesparta/fullcycle-clean-architecture/pkg/events"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfgs, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(
		cfgs.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfgs.DBUser, cfgs.DBPassword, cfgs.DBHost, cfgs.DBPort, cfgs.DBName),
	)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(cfgs.RabbitMqURI)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db, eventDispatcher)

	newWebServer := webserver.NewWebServer(cfgs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	newWebServer.AddHandler(http.MethodPost, "/order", webOrderHandler.Create)
	newWebServer.AddHandler(http.MethodGet, "/order", webOrderHandler.List)
	fmt.Println("Starting web server on port", cfgs.WebServerPort)
	go newWebServer.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", cfgs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfgs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphqlhandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &gql.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", cfgs.GraphQLServerPort)
	http.ListenAndServe(":"+cfgs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(rabbitmqURI string) *amqp.Channel {
	conn, err := amqp.Dial(rabbitmqURI)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
