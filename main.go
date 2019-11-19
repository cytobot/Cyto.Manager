package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

// VERSION is the application version
const VERSION = "0.1.0"

type managerState struct {
	nats       *NatsClient
	data       *DatabaseClient
	grpcServer *grpc.Server
}

func main() {
	dbClient := getDatabaseClient()

	manager := &managerState{
		nats:       getNatsClient(),
		data:       dbClient,
		grpcServer: NewManagerServer(dbClient.CommandRepository),
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go manager.grpcServer.Serve(lis)

	log.Println("Started successfully")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

out:
	for {
		select {
		case <-c:
			log.Println("Shutting down...")
			break out
		}
	}
}

func getNatsClient() *NatsClient {
	natsEndpoint := os.Getenv("NatsEndpoint")

	if natsEndpoint == "" {
		panic("No nats endpoint provided.")
	}

	client, err := NewNatsClient(natsEndpoint)
	if err != nil {
		panic(fmt.Sprintf("[NATS error] %s", err))
	}

	log.Println("Connected to NATS")

	return client
}

func getDatabaseClient() *DatabaseClient {
	connectionString := os.Getenv("DBConnectionString")

	if connectionString == "" {
		panic("No database connection string provided.")
	}

	client, err := NewDatabaseClient(connectionString)
	if err != nil {
		panic(fmt.Sprintf("[Database error] %s", err))
	}

	log.Println("Connected to database")

	return client
}
