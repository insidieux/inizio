package plugins

import (
	"context"
	"testing"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestCapitaliseFirst(t *testing.T) {
	suite.Run(t, new(capitaliseFirstTestSuite))
}

func TestRegistry(t *testing.T) {
	suite.Run(t, new(registryTestSuite))
}

// --- Suites ---

type capitaliseFirstTestSuite struct {
	suite.Suite
}

func (s *capitaliseFirstTestSuite) TestCapitaliseFirst() {
	s.Equal(``, capitaliseFirst(``))
	s.Equal(`First second`, capitaliseFirst(`first second`))
	s.Equal(`First Second`, capitaliseFirst(`first Second`))
	s.Equal(`First Second`, capitaliseFirst(`First Second`))
}

type registryTestSuite struct {
	suite.Suite
	context context.Context
	options generator.RunOptions
	values  generator.RunValues
}

func (s *registryTestSuite) SetupTest() {
	s.context = context.Background()
	s.options = generator.RunOptions{}
	s.values = generator.RunValues{}
}

func (s *registryTestSuite) TestHasPlugins() {
	registry := NewRegistry(nil, nil, false)
	s.False(registry.HasPlugins())
	mockedClient := new(mockClientInterface)
	registry = NewRegistry([]ClientInterface{mockedClient}, nil, false)
	s.True(registry.HasPlugins())
}

func (s *registryTestSuite) TestProcessError() {
	mockedClient := new(mockClientInterface)
	mockedClient.On(`Name`).Return(`name`)
	mockedClient.On(`Run`, s.context, s.options, s.values).Return(nil, errors.New(`expected error`))

	registry := NewRegistry([]ClientInterface{mockedClient}, logrus.New(), true)
	result, err := registry.Process(s.context, s.options, s.values)
	s.Error(err)
	s.EqualError(err, `plugin "name" error: expected error`)
	s.Nil(result)
}

func (s *registryTestSuite) TestProcess() {
	mockedFirstClient := new(mockClientInterface)
	mockedFirstClient.On(`Name`).Return(`first`)
	mockedFirstClient.On(`Run`, s.context, s.options, s.values).Return(nil, errors.New(`expected error`))

	mockedSecondClient := new(mockClientInterface)
	mockedSecondClient.On(`Name`).Return(`second`)
	mockedSecondClient.On(`Run`, s.context, s.options, s.values).Return(generator.RunResult{`/one`, `/two`}, nil)

	registry := NewRegistry([]ClientInterface{mockedFirstClient, mockedSecondClient}, logrus.New(), false)
	result, err := registry.Process(s.context, s.options, s.values)
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 2)
	s.Equal(generator.RunResult{`"second" generated: /one`, `"second" generated: /two`}, result)
}
