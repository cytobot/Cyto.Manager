package main

import (
	"encoding/json"
	"log"

	nats "github.com/nats-io/nats.go"
)

type NatsClient struct {
	client *nats.Conn
}

func NewNatsClient(endpoint string) *NatsClient {
	client, err := nats.Connect(endpoint)
	defer client.Drain()

	if err != nil {
		panic("Unable to connect to nats")
	}

	return &NatsClient{
		client: client,
	}
}

func (c *NatsClient) Publish(subject string, msg interface{}) error {
	bytes, err := json.Marshal(msg)

	if err != nil {
		log.Printf("Failed to marshal message: %s", err)
		return err
	}

	err = c.client.Publish(subject, bytes)

	if err != nil {
		log.Printf("Failed to publish message: %s", err)
		return err
	}

	return nil
}
