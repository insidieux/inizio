package plugin

import (
	"testing"

	pluginSDK "github.com/hashicorp/go-plugin"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// --- Tests ---

func TestGRPCPlugin(t *testing.T) {
	suite.Run(t, new(gRPCPluginTestSuite))
}

// --- Suites ---

type gRPCPluginTestSuite struct {
	suite.Suite
}

func (s *gRPCPluginTestSuite) TestGRPCServer() {
	plugin := new(gRPCPlugin)
	plugin.generator = new(mockGenerator)

	s.NoError(plugin.GRPCServer(nil, pluginSDK.DefaultGRPCServer(nil)))
}

func (s *gRPCPluginTestSuite) TestGRPCClient() {
	plugin := new(gRPCPlugin)
	plugin.generator = new(mockGenerator)

	client, err := plugin.GRPCClient(nil, nil, new(grpc.ClientConn))
	s.NoError(err)
	s.Implements(new(generator.Generator), client)
}
