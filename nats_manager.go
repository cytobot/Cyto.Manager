package main

import (
	"encoding/json"
	"strconv"
	"time"

	cytonats "github.com/cytobot/messaging/nats"
	pbd "github.com/cytobot/messaging/transport/discord"
	pbm "github.com/cytobot/messaging/transport/manager"
)

type NatsManager struct {
	client       *cytonats.NatsClient
	listenerChan chan int32
}

func NewNatsManager(endpoint string) (*NatsManager, error) {
	client, err := cytonats.NewNatsClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &NatsManager{
		client: client,
	}, nil
}

// TODO: Keep track of individual shard state for cluster reporting.
func (m *NatsManager) StartHealthCheckListener() error {
	listenerHealth, err := m.client.ChanSubscribe("listener_health")
	if err != nil {
		return err
	}

	m.listenerChan = make(chan int32)
	go func() {
		for {
			select {
			case msg := <-listenerHealth:
				healthMsg := &pbd.HealthCheckStatus{}
				json.Unmarshal(msg.Data, healthMsg)

				//log.Printf("[Discord Health] %+v", healthMsg)
			case <-m.listenerChan:
				return
			}
		}
	}()

	return nil
}

func (m *NatsManager) NotifyCommandConfigChange(updatedCommandDefinitions []*CommandDefinition) error {
	protoDefinitions := convertToProtoCommandDefinitions(updatedCommandDefinitions)

	natsUpdate := &pbm.UpdatedCommandConfigurations{
		Timestamp:          mapToProtoTimestamp(time.Now().UTC()),
		CommandDefinitions: protoDefinitions,
	}

	return m.client.Publish("command-update", natsUpdate)
}

func (m *NatsManager) Shutdown() {
	if m.listenerChan != nil {
		m.listenerChan <- 0
	}
	m.client.Shutdown()
}

func convertToProtoCommandDefinitions(commandDefinitions []*CommandDefinition) []*pbm.CommandDefinition {
	protoDefinitions := make([]*pbm.CommandDefinition, 0)

	for _, def := range commandDefinitions {
		newDef := &pbm.CommandDefinition{
			CommandID:            def.CommandID,
			Enabled:              def.Enabled,
			Triggers:             def.Triggers,
			PermissionLevel:      mapToNatsPermissionLevel(def.PermissionLevel),
			ParameterDefinitions: mapToNatsParameterDefinition(def.ParameterDefinitions),
			LastModifiedUserID:   def.LastModifiedUserID,
			LastModifiedDateUtc:  mapToProtoTimestamp(def.LastModifiedDateUtc),
		}
		protoDefinitions = append(protoDefinitions, newDef)
	}

	return protoDefinitions
}

func mapToNatsPermissionLevel(permissionLevel string) pbm.CommandDefinition_PermissionLevel {
	protoEnumValue := pbm.CommandDefinition_PermissionLevel_value[permissionLevel]
	return pbm.CommandDefinition_PermissionLevel(protoEnumValue)
}

func mapToNatsParameterDefinition(commandParameterDefinitions []CommandParameterDefinition) []*pbm.CommandParameterDefinition {
	protoParameterDefinitions := make([]*pbm.CommandParameterDefinition, 0)
	for _, p := range commandParameterDefinitions {
		protoParameterDefinitions = append(protoParameterDefinitions, &pbm.CommandParameterDefinition{
			Name:     p.Name,
			Pattern:  p.Pattern,
			Optional: strconv.FormatBool(p.Optional),
		})
	}
	return protoParameterDefinitions
}
