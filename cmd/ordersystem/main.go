package main

import (
	"database/sql"
	"fmt"
	"github.com/andre2ar/orders/internal/infra/web"
	"net"
	"net/http"

	graphqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/andre2ar/orders/configs"
	"github.com/andre2ar/orders/internal/event/handler"
	"github.com/andre2ar/orders/internal/infra/graph"
	"github.com/andre2ar/orders/internal/infra/grpc/pb"
	"github.com/andre2ar/orders/internal/infra/grpc/service"
	"github.com/andre2ar/orders/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	appConfigs, err := configs.LoadConfig("./cmd/ordersystem")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(appConfigs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", appConfigs.DBUser, appConfigs.DBPassword, appConfigs.DBHost, appConfigs.DBPort, appConfigs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(appConfigs.RabbitMqUrl)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	restServer := web.NewWebServer(appConfigs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)

	restServer.AddHandler(web.POST, "/order", webOrderHandler.Create)
	restServer.AddHandler(web.GET, "/order", webOrderHandler.List)

	fmt.Println("Starting web server on port", appConfigs.WebServerPort)
	go restServer.Start()

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrdersUseCase(db)

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", appConfigs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfigs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	graphQlHandler := graphqlHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", graphQlHandler)

	fmt.Println("Starting GraphQL server on port", appConfigs.GraphQLServerPort)
	http.ListenAndServe(":"+appConfigs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(url string) *amqp.Channel {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
