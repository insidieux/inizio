package plugins

import (
	"context"
	"testing"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestClient(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}

// --- Suites ---

type clientTestSuite struct {
	suite.Suite
	ctx     context.Context
	options generator.RunOptions
	values  generator.RunValues
}

func (s *clientTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.options = generator.RunOptions{}
	s.values = generator.RunValues{}
}

func (s *clientTestSuite) TestClientName() {
	c := new(client)
	c.name = `name`

	s.Equal(`name`, c.Name())
}

func (s *clientTestSuite) TestRunError() {
	mockedClient := new(mockPluginClientInterface)
	mockedClient.On(`Dispense`).Return(nil, errors.New(`expected error`))

	c := new(client)
	c.downstream = mockedClient

	result, err := c.Run(s.ctx, s.options, s.values)
	s.Error(err)
	s.Nil(result)
}

func (s *clientTestSuite) TestRunErrorGenerator() {
	mockedGenerator := new(mockGenerator)
	mockedGenerator.On(`Run`, s.ctx, s.options, s.values).Return(nil, errors.New(`expected error`))

	mockedClient := new(mockPluginClientInterface)
	mockedClient.On(`Dispense`).Return(mockedGenerator, nil)

	c := new(client)
	c.downstream = mockedClient

	result, err := c.Run(s.ctx, s.options, s.values)
	s.Error(err)
	s.Nil(result)
}

func (s *clientTestSuite) TestRun() {
	mockedGenerator := new(mockGenerator)
	mockedGenerator.On(`Run`, s.ctx, s.options, s.values).Return(generator.RunResult{`one`, `two`}, nil)

	mockedClient := new(mockPluginClientInterface)
	mockedClient.On(`Dispense`).Return(mockedGenerator, nil)

	c := new(client)
	c.downstream = mockedClient

	result, err := c.Run(s.ctx, s.options, s.values)
	s.NoError(err)
	s.Equal(generator.RunResult{`one`, `two`}, result)
}
