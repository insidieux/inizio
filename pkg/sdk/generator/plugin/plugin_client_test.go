package plugin

import (
	"context"
	"errors"
	"os/exec"
	"testing"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestPluginClient(t *testing.T) {
	suite.Run(t, new(pluginClientTestSuite))
}

// --- Suites ---

type pluginClientTestSuite struct {
	suite.Suite
}

func (s *pluginClientTestSuite) TestNewClient() {
	c := NewClient(exec.CommandContext(context.Background(), `name`))
	s.Implements(new(ClientInterface), c)
}

func (s *pluginClientTestSuite) TestClientDispenseErrorConnectRPC() {
	mockedSDKClient := new(mockSdkClient)
	mockedSDKClient.On(`Client`).Return(nil, errors.New(`expected error`))

	c := new(client)
	c.downstream = mockedSDKClient

	gen, err := c.Dispense()
	s.Error(err)
	s.EqualError(err, `failed to connect via RPC: expected error`)
	s.Nil(gen)
}

func (s *pluginClientTestSuite) TestClientDispenseErrorDispense() {
	mockedClientProtocol := new(mockClientProtocol)
	mockedClientProtocol.On(`Dispense`, pluginID).Return(nil, errors.New(`expected error`))

	mockedSDKClient := new(mockSdkClient)
	mockedSDKClient.On(`Client`).Return(mockedClientProtocol, nil)

	c := new(client)
	c.downstream = mockedSDKClient

	gen, err := c.Dispense()
	s.Error(err)
	s.EqualError(err, `failed to dispense generator plugin: expected error`)
	s.Nil(gen)
}

func (s *pluginClientTestSuite) TestClientDispense() {
	grpcClient := new(gRPCClient)

	mockedClientProtocol := new(mockClientProtocol)
	mockedClientProtocol.On(`Dispense`, pluginID).Return(grpcClient, nil)

	mockedSDKClient := new(mockSdkClient)
	mockedSDKClient.On(`Client`).Return(mockedClientProtocol, nil)

	c := new(client)
	c.downstream = mockedSDKClient

	gen, err := c.Dispense()
	s.NoError(err)
	s.NotNil(gen)
	s.Implements(new(generator.Generator), gen)
}
