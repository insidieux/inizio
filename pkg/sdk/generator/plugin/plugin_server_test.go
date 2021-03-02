package plugin

import (
	"context"
	"sync"
	"testing"

	pluginSDK "github.com/hashicorp/go-plugin"

	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestPluginServer(t *testing.T) {
	suite.Run(t, new(pluginServerTestSuite))
}

// --- Suites ---

type pluginServerTestSuite struct {
	suite.Suite
}

func (s *pluginServerTestSuite) TestWithContext() {
	ctx := context.Background()
	option := WithContext(ctx)
	config := new(pluginSDK.ServeConfig)
	option(config)
	s.Equal(ctx, config.Test.Context)
}

func (s *pluginServerTestSuite) TestServe() {
	ctx, cancel := context.WithCancel(context.Background())
	mockedGenerator := new(mockGenerator)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		Serve(mockedGenerator, WithContext(ctx))
	}()
	cancel()
	wg.Wait()
}
