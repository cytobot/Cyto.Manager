package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

type ManagerState struct {
	nats          *NatsManager
	data          *DatabaseClient
	grpcServer    *grpc.Server
	healthMonitor *healthMonitor
}

func main() {
	dbClient := getDatabaseClient()
	healthMonitor := NewHealthMonitor()

	manager := &ManagerState{
		data:          dbClient,
		healthMonitor: healthMonitor,
		grpcServer:    NewRpcServer(dbClient.CommandRepository, healthMonitor),
	}

	manager.nats = getNatsManager(manager)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go manager.grpcServer.Serve(lis)
	go manager.nats.StartHealthCheckListener()

	log.Println("Started successfully")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

out:
	for {
		select {
		case <-c:
			manager.nats.Shutdown()
			log.Println("Shutting down...")
			break out
		}
	}
}

func getNatsManager(s *ManagerState) *NatsManager {
	natsEndpoint := os.Getenv("NatsEndpoint")

	if natsEndpoint == "" {
		panic("No nats endpoint provided.")
	}

	manager, err := NewNatsManager(natsEndpoint, s)
	if err != nil {
		panic(fmt.Sprintf("[NATS error] %s", err))
	}

	log.Println("Connected to NATS")

	return manager
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
