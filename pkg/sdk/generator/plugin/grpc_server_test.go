package plugin

import (
	"context"
	"errors"
	"testing"

	pluginSDK "github.com/hashicorp/go-plugin"

	"github.com/insidieux/inizio/pkg/api/protobuf"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestGRPCServer(t *testing.T) {
	suite.Run(t, new(gRPCServerTestSuite))
}

// --- Suites ---

type gRPCServerTestSuite struct {
	suite.Suite
	request *protobuf.Run_Request
}

func (s *gRPCServerTestSuite) SetupTest() {
	s.request = &protobuf.Run_Request{
		Options: &protobuf.Run_Request_Options{
			Override:         false,
			WorkingDirectory: `/dir`,
		},
		Values: &protobuf.Run_Request_Values{
			Application: &protobuf.Run_Request_Values_Application{},
			Golang:      &protobuf.Run_Request_Values_Golang{},
		},
	}
}

func (s *gRPCServerTestSuite) TestGRPCServer() {
	plugin := new(gRPCPlugin)
	s.NoError(plugin.GRPCServer(nil, pluginSDK.DefaultGRPCServer(nil)))
}

func (s *gRPCServerTestSuite) TestRunDownstreamError() {
	ctx := context.Background()

	downstream := new(mockGenerator)
	downstream.On(`Run`, ctx, mock.Anything, mock.Anything).
		Return(nil, errors.New(`expected error`))

	c := new(gRPCServer)
	c.downstream = downstream

	result, err := c.Run(ctx, s.request)
	s.Error(err)
	s.Nil(result)
}

func (s *gRPCServerTestSuite) TestRun() {
	ctx := context.Background()

	downstream := new(mockGenerator)
	downstream.On(`Run`, ctx, mock.Anything, mock.Anything).
		Return(generator.RunResult([]string{`/some/path/first`, `/some/path/second`}), nil)

	c := new(gRPCServer)
	c.downstream = downstream

	result, err := c.Run(ctx, s.request)
	s.NoError(err)
	s.NotEmpty(result.Generated)
	s.Len(result.Generated, 2)
}
