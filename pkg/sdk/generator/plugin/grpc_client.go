package plugin

import (
	"context"

	"github.com/insidieux/inizio/pkg/api/protobuf"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
)

type (
	gRPCClient struct {
		downstream protobuf.GeneratorClient
	}
)

var (
	_ generator.Generator = &gRPCClient{}
)

func (s *gRPCClient) Run(ctx context.Context, options generator.RunOptions, values generator.RunValues) (generator.RunResult, error) {
	result, err := s.downstream.Run(ctx, generatorRunArgumentsToProto(options, values))
	if err != nil {
		return nil, errors.Wrap(err, `grpc request failed`)
	}
	return protoRunResponseToResult(result), nil
}
