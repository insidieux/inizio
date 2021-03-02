package plugin

import (
	"os"
	"os/exec"

	pluginSDK "github.com/hashicorp/go-plugin"

	"github.com/hashicorp/go-hclog"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
)

type (
	// ClientInterface hold information about further GRPC connection and knows how to Dispense plugin.
	ClientInterface interface {
		Dispense() (generator.Generator, error)
	}
	sdkClient interface {
		Client() (pluginSDK.ClientProtocol, error)
	}
	client struct {
		downstream sdkClient
	}
)

// NewClient make GRPC client to plugin implementation
func NewClient(cmd *exec.Cmd) ClientInterface {
	return &client{
		downstream: pluginSDK.NewClient(&pluginSDK.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins: pluginSDK.PluginSet{
				pluginID: &gRPCPlugin{},
			},
			AllowedProtocols: []pluginSDK.Protocol{
				pluginSDK.ProtocolGRPC,
			},
			// watch documentation for pluginSDK.ClientConfig.Managed - it used to call Kill by system
			Managed: true,
			Cmd:     cmd,
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:   `plugin`,
				Level:  hclog.NoLevel,
				Output: os.Stdout,
			}),
		}),
	}
}

// Dispense create GRPC client, make dispense call to plugin and return gRPCClient/generator.Generator implementation
func (c *client) Dispense() (generator.Generator, error) {
	// Connect via RPC
	rpc, err := c.downstream.Client()
	if err != nil {
		return nil, errors.Wrap(err, `failed to connect via RPC`)
	}

	// Request generator plugin
	dispense, err := rpc.Dispense(pluginID)
	if err != nil {
		return nil, errors.Wrap(err, `failed to dispense generator plugin`)
	}

	// Create plugin instance
	return dispense.(*gRPCClient), nil
}
