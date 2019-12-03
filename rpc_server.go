package main

import (
	"context"
	"time"

	pb "github.com/cytobot/rpc/manager"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

type rpcServer struct {
	pb.UnimplementedManagerServer
	repository    *CommandRepository
	healthMonitor *healthMonitor
}

func NewRpcServer(repository *CommandRepository, healthMonitor *healthMonitor) *grpc.Server {
	rpcServer := &rpcServer{
		repository:    repository,
		healthMonitor: healthMonitor,
	}
	server := grpc.NewServer()
	pb.RegisterManagerServer(server, rpcServer)
	return server
}

func (s *rpcServer) GetCommandDefinitions(ctx context.Context, req *empty.Empty) (*pb.CommandDefinitionList, error) {
	commandDefintions, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	protoCommandDefinition := make([]*pb.CommandDefinition, 0)

	for _, def := range commandDefintions {

		pDef := &pb.CommandDefinition{
			CommandID:            def.CommandID,
			Enabled:              def.Enabled,
			Unlisted:             def.Unlisted,
			Description:          def.Description,
			Triggers:             def.Triggers,
			PermissionLevel:      mapToProtoPermissionLevel(def.PermissionLevel),
			ParameterDefinitions: mapToProtoParameterDefinition(def.ParameterDefinitions),
			LastModifiedUserID:   def.LastModifiedUserID,
			LastModifiedDateUtc:  MapToProtoTimestamp(def.LastModifiedDateUtc),
		}
		protoCommandDefinition = append(protoCommandDefinition, pDef)
	}

	return &pb.CommandDefinitionList{
		CommandDefinitions: protoCommandDefinition,
	}, nil
}

func (s *rpcServer) GetGuildCommandConfigurations(ctx context.Context, req *pb.GuildQuery) (*pb.GuildCommandConfigurationList, error) {
	return nil, nil
}

func (s *rpcServer) SetGuildCommandConfiguration(ctx context.Context, req *pb.GuildCommandConfiguration) (*pb.GuildCommandConfiguration, error) {
	return nil, nil
}

func (s *rpcServer) GetWorkerHealthChecks(ctx context.Context, req *empty.Empty) (*pb.HealthCheckStatusList, error) {
	statusList := s.healthMonitor.GetAllWorkerStatus()

	return &pb.HealthCheckStatusList{
		HealthChecks: statusList,
	}, nil
}

func (s *rpcServer) GetListenerHealthChecks(ctx context.Context, req *empty.Empty) (*pb.HealthCheckStatusList, error) {
	statusList := s.healthMonitor.GetAllListenerStatus()

	return &pb.HealthCheckStatusList{
		HealthChecks: statusList,
	}, nil
}

func MapToProtoTimestamp(timeValue time.Time) *timestamp.Timestamp {
	protoTimestamp, _ := ptypes.TimestampProto(timeValue)
	return protoTimestamp
}

func MapFromProtoTimestamp(ts *timestamp.Timestamp) time.Time {
	timeValue, _ := ptypes.Timestamp(ts)
	return timeValue
}

func mapToProtoPermissionLevel(permissionLevel string) pb.CommandDefinition_PermissionLevel {
	protoEnumValue := pb.CommandDefinition_PermissionLevel_value[permissionLevel]
	return pb.CommandDefinition_PermissionLevel(protoEnumValue)
}

func mapToProtoParameterDefinition(commandParameterDefinitions []CommandParameterDefinition) []*pb.CommandParameterDefinition {
	protoParameterDefinitions := make([]*pb.CommandParameterDefinition, 0)
	for _, p := range commandParameterDefinitions {
		protoParameterDefinitions = append(protoParameterDefinitions, &pb.CommandParameterDefinition{
			Name:     p.Name,
			Pattern:  p.Pattern,
			Optional: p.Optional,
		})
	}
	return protoParameterDefinitions
}
