package main

import (
	"time"
)

type cytoCommandConfigUpdate struct {
	Timestamp          time.Time
	CommandDefinitions []*CommandDefinition
}

func (c *NatsClient) notifyCommandConfigChange(updatedCommandDefinitions []*CommandDefinition) error {
	natsUpdate := &cytoCommandConfigUpdate{
		Timestamp:          time.Now(),
		CommandDefinitions: updatedCommandDefinitions,
	}

	return c.Publish("command-update", natsUpdate)
}
