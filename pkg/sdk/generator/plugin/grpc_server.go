package plugin

import (
	"context"

	"github.com/insidieux/inizio/pkg/api/protobuf"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
)

type (
	gRPCServer struct {
		protobuf.UnimplementedGeneratorServer
		downstream generator.Generator
	}
)

var (
	_ protobuf.GeneratorServer = &gRPCServer{}
)

func (s *gRPCServer) Run(ctx context.Context, request *protobuf.Run_Request) (*protobuf.Run_Response, error) {
	result, err := s.downstream.Run(ctx, protoRunRequestToOptions(request), protoRunRequestToValues(request))
	if err != nil {
		return nil, errors.Wrap(err, `failed to run downstream GRPC plugin request`)
	}
	return generatorRunResultToProto(result), nil
}
