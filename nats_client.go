package main

import (
	"encoding/json"
	"log"

	nats "github.com/nats-io/nats.go"
)

type NatsClient struct {
	client *nats.Conn
}

func NewNatsClient(endpoint string) (*NatsClient, error) {
	client, err := nats.Connect(endpoint)
	if err != nil {
		return nil, err
	}

	defer client.Drain()

	return &NatsClient{
		client: client,
	}, nil
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
