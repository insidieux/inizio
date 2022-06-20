package plugins

import (
	"context"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/insidieux/inizio/pkg/sdk/generator/plugin"
	"github.com/pkg/errors"
)

type (
	// ClientInterface is common interface, holding information about GRPC connection to plugin endpoint.
	ClientInterface interface {
		generator.Generator
		Name() string
	}
	client struct {
		name       string
		downstream plugin.ClientInterface
	}
)

var (
	_ ClientInterface = &client{}
)

// Name return shortcut name of plugin, without prefix and version. Used for config and bootstrap.
func (c *client) Name() string {
	return c.name
}

// Run create GRPC client, make dispense call to plugin via downstream and call downstream run method.
func (c *client) Run(ctx context.Context, options generator.RunOptions, values generator.RunValues) (generator.RunResult, error) {
	gen, err := c.downstream.Dispense()
	if err != nil {
		return nil, errors.Wrap(err, `bootstrap plugin failed`)
	}
	result, err := gen.Run(ctx, options, values)
	if err != nil {
		return nil, errors.Wrap(err, `generation failed`)
	}
	return result, nil
}
