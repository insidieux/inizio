package plugin

import (
	"context"
	"os"

	pluginSDK "github.com/hashicorp/go-plugin"

	"github.com/hashicorp/go-hclog"
	"github.com/insidieux/inizio/pkg/sdk/generator"
)

type (
	// ServeOption allows setting optional pluginSDK.ServeConfig values.
	ServeOption func(*pluginSDK.ServeConfig)
)

// WithContext set context.Context to pluginSDK.ServeConfig and can be used for testing cased.
func WithContext(ctx context.Context) ServeOption {
	return func(config *pluginSDK.ServeConfig) {
		if config.Test == nil {
			config.Test = new(pluginSDK.ServeTestConfig)
		}
		config.Test.Context = ctx
	}
}

// Serve run plugin implementation as GRPC server.
func Serve(implementation generator.Generator, options ...ServeOption) {
	config := &pluginSDK.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: pluginSDK.PluginSet{
			pluginID: &gRPCPlugin{
				generator: implementation,
			},
		},
		GRPCServer: pluginSDK.DefaultGRPCServer,
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   `plugin`,
			Level:  hclog.NoLevel,
			Output: os.Stdout,
		}),
	}
	for _, option := range options {
		option(config)
	}
	pluginSDK.Serve(config)
}
