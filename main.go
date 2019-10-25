package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	cydata "./lib/data"
	"google.golang.org/grpc"
)

// VERSION is the application version
const VERSION = "0.1.0"

type managerState struct {
	nats       *NatsClient
	data       *cydata.DatabaseClient
	grpcServer *grpc.Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dbClient := getDatabaseClient()

	manager := &managerState{
		nats:       getNatsClient(),
		data:       dbClient,
		grpcServer: NewManagerServer(dbClient.CommandRepository),
	}

	manager.grpcServer.Serve(lis)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

out:
	for {
		select {
		case <-c:
			break out
		}
	}
}

func getNatsClient() *NatsClient {
	natsEndpoint := os.Getenv("NatsEndpoint")

	if natsEndpoint == "" {
		panic("No nats endpoint provided.")
	}

	return NewNatsClient(natsEndpoint)
}

func getDatabaseClient() *cydata.DatabaseClient {
	connectionString := os.Getenv("DBConnectionString")

	if connectionString == "" {
		panic("No database connection string provided.")
	}

	client, err := cydata.NewDatabaseClient("sqlite3", connectionString)
	if err != nil {
		panic(fmt.Sprintf("[Database error] %s", err))
	}

	return client
}
