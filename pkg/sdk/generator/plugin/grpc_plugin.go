package plugin

import (
	"context"

	pluginSDK "github.com/hashicorp/go-plugin"

	"github.com/insidieux/inizio/pkg/api/protobuf"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"google.golang.org/grpc"
)

const (
	pluginVersion = 1
	pluginID      = `generator`

	magicCookieKey   = `INIZIO_GENERATOR_PLUGIN`
	magicCookieValue = `c7f8f3fff5841c03cc871e3606c9cfdc` // DO NOT CHANGE THIS VALUE
)

type (
	gRPCPlugin struct {
		pluginSDK.NetRPCUnsupportedPlugin
		generator generator.Generator
	}
)

var (
	handshakeConfig = pluginSDK.HandshakeConfig{
		ProtocolVersion:  pluginVersion,
		MagicCookieKey:   magicCookieKey,
		MagicCookieValue: magicCookieValue,
	}
	_ pluginSDK.Plugin     = &gRPCPlugin{}
	_ pluginSDK.GRPCPlugin = &gRPCPlugin{}
)

func (p *gRPCPlugin) GRPCServer(_ *pluginSDK.GRPCBroker, s *grpc.Server) error {
	protobuf.RegisterGeneratorServer(s, &gRPCServer{
		downstream: p.generator,
	})
	return nil
}

func (p *gRPCPlugin) GRPCClient(_ context.Context, _ *pluginSDK.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &gRPCClient{downstream: protobuf.NewGeneratorClient(c)}, nil
}
