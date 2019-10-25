package main

import (
	"time"

	"./lib/data"
)

type cytoCommandConfigUpdate struct {
	Timestamp          time.Time
	CommandDefinitions []data.CommandDefinition
}

func (c *NatsClient) notifyCommandConfigChange(updatedCommandDefinitions []data.CommandDefinition) error {
	natsUpdate := &cytoCommandConfigUpdate{
		Timestamp:          time.Now(),
		CommandDefinitions: updatedCommandDefinitions,
	}

	return c.Publish("command-update", natsUpdate)
}
