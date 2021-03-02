package plugin

import (
	"context"
	"testing"

	"github.com/insidieux/inizio/pkg/api/protobuf"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestGRPCClient(t *testing.T) {
	suite.Run(t, new(gRPCClientTestSuite))
}

// --- Suites ---

type gRPCClientTestSuite struct {
	suite.Suite
}

func (s *gRPCClientTestSuite) TestRunDownstreamError() {
	ctx := context.Background()

	downstream := new(mockGeneratorClient)
	downstream.On(`Run`, ctx, mock.Anything, mock.Anything).Return(nil, errors.New(`expected error`))

	c := new(gRPCClient)
	c.downstream = downstream

	result, err := c.Run(ctx, generator.RunOptions{}, generator.RunValues{})
	s.Error(err)
	s.Nil(result)
}

func (s *gRPCClientTestSuite) TestRun() {
	ctx := context.Background()

	downstream := new(mockGeneratorClient)
	downstream.On(`Run`, ctx, mock.Anything, mock.Anything).Return(&protobuf.Run_Response{
		Generated: []string{`/some/path/first`, `/some/path/second`},
	}, nil)

	c := new(gRPCClient)
	c.downstream = downstream

	result, err := c.Run(ctx, generator.RunOptions{}, generator.RunValues{})
	s.NoError(err)
	s.NotEmpty(result)
	s.Len(result, 2)
}
