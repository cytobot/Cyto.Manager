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
	state        *ManagerState
	shutdownChan chan int32
}

func NewNatsManager(endpoint string, state *ManagerState) (*NatsManager, error) {
	client, err := cytonats.NewNatsClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &NatsManager{
		client:       client,
		state:        state,
		shutdownChan: make(chan int32),
	}, nil
}

// TODO: Keep track of individual shard state for cluster reporting.
func (m *NatsManager) StartHealthCheckListener() error {
	listenerHealth, err := m.client.ChanSubscribe("listener_health")
	if err != nil {
		return err
	}

	workerHealth, err := m.client.ChanSubscribe("worker_health")
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-listenerHealth:
				healthMsg := &pbd.HealthCheckStatus{}
				json.Unmarshal(msg.Data, healthMsg)

				m.state.healthMonitor.UpdateListenerHealthStatus(healthMsg)
			case <-m.shutdownChan:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case msg := <-workerHealth:
				healthMsg := &pbd.HealthCheckStatus{}
				json.Unmarshal(msg.Data, healthMsg)

				m.state.healthMonitor.UpdateWorkerHealthStatus(healthMsg)
			case <-m.shutdownChan:
				return
			}
		}
	}()

	return nil
}

func (m *NatsManager) NotifyCommandConfigChange(updatedCommandDefinitions []*CommandDefinition) error {
	protoDefinitions := convertToProtoCommandDefinitions(updatedCommandDefinitions)

	natsUpdate := &pbm.UpdatedCommandConfigurations{
		Timestamp:          MapToProtoTimestamp(time.Now().UTC()),
		CommandDefinitions: protoDefinitions,
	}

	return m.client.Publish("command-update", natsUpdate)
}

func (m *NatsManager) Shutdown() {
	m.shutdownChan <- 0
	m.client.Shutdown()
}

func convertToProtoCommandDefinitions(commandDefinitions []*CommandDefinition) []*pbm.CommandDefinition {
	protoDefinitions := make([]*pbm.CommandDefinition, 0)

	for _, def := range commandDefinitions {
		newDef := &pbm.CommandDefinition{
			CommandID:            def.CommandID,
			Enabled:              def.Enabled,
			Unlisted:             def.Unlisted,
			Description:          def.Description,
			Triggers:             def.Triggers,
			PermissionLevel:      mapToNatsPermissionLevel(def.PermissionLevel),
			ParameterDefinitions: mapToNatsParameterDefinition(def.ParameterDefinitions),
			LastModifiedUserID:   def.LastModifiedUserID,
			LastModifiedDateUtc:  MapToProtoTimestamp(def.LastModifiedDateUtc),
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
