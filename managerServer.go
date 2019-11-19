package main

import (
	"context"
	"strconv"
	"time"

	pb "github.com/cytobot/rpc/manager"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

func NewManagerServer(repository *CommandRepository) *grpc.Server {
	managerServer := &managerServer{
		repository: repository,
	}
	server := grpc.NewServer()
	pb.RegisterManagerServer(server, managerServer)
	return server
}

type managerServer struct {
	pb.UnimplementedManagerServer
	repository *CommandRepository
}

func (s *managerServer) GetCommandDefinitions(ctx context.Context, req *empty.Empty) (*pb.CommandDefinitionList, error) {
	commandDefintions, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	protoCommandDefinition := make([]*pb.CommandDefinition, 0)

	for _, def := range commandDefintions {

		pDef := &pb.CommandDefinition{
			CommandID:            def.CommandID,
			Enabled:              def.Enabled,
			Triggers:             def.Triggers,
			PermissionLevel:      mapToProtoPermissionLevel(def.PermissionLevel),
			ParameterDefinitions: mapToProtoParameterDefinition(def.ParameterDefinitions),
			LastModifiedUserID:   def.LastModifiedUserID,
			LastModifiedDateUtc:  mapToProtoTimestamp(def.LastModifiedDateUtc),
		}
		protoCommandDefinition = append(protoCommandDefinition, pDef)
	}

	return &pb.CommandDefinitionList{
		CommandDefinitions: protoCommandDefinition,
	}, nil
}

func (s *managerServer) GetGuildCommandConfigurations(ctx context.Context, req *pb.GuildQuery) (*pb.GuildCommandConfigurationList, error) {
	return nil, nil
}

func (s *managerServer) SetGuildCommandConfiguration(ctx context.Context, req *pb.GuildCommandConfiguration) (*pb.GuildCommandConfiguration, error) {
	return nil, nil
}

func mapToProtoTimestamp(timeValue time.Time) *timestamp.Timestamp {
	protoTimestamp, _ := ptypes.TimestampProto(timeValue)
	return protoTimestamp
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
			Optional: strconv.FormatBool(p.Optional),
		})
	}
	return protoParameterDefinitions
}
