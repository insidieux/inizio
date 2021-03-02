package plugins

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/drone/envsubst"
	"github.com/pkg/errors"
)

type (
	// PluginConfigInterface determine how to make exec.Command from passed configuration .
	PluginConfigInterface interface {
		Command(context.Context, string, string) (*exec.Cmd, error)
	}
	// ConfigInterface determine how to find PluginConfigInterface by name.
	ConfigInterface interface {
		Lookup(string) (PluginConfigInterface, error)
	}
	// Config is ConfigInterface implementation.
	Config []*PluginConfig
	// PluginConfig is PluginConfigInterface implementation.
	PluginConfig struct {
		Name string   `json:"name" yaml:"name"`
		Args []string `json:"args" yaml:"args"`
		Env  []struct {
			Name  string `json:"name" yaml:"name"`
			Value string `json:"value" yaml:"value"`
		} `json:"env" yaml:"env"`
		Flags []struct {
			Name  string `json:"name" yaml:"name"`
			Value string `json:"value" yaml:"value"`
		} `json:"flags" yaml:"flags"`
	}
)

func newPluginConfig(name string) PluginConfigInterface {
	cfg := new(PluginConfig)
	cfg.Name = name
	return cfg
}

// Lookup iterate over slice of PluginConfigInterface and try to find suitable by name.
func (c Config) Lookup(name string) (PluginConfigInterface, error) {
	for _, item := range c {
		if item.Name != name {
			continue
		}
		return item, nil
	}
	return nil, errors.Errorf(`plugin "%s" is not defined in config`, name)
}

// Command collect information from passed config, modify data and create exec.Command.
func (pc *PluginConfig) Command(ctx context.Context, path string, executable string) (*exec.Cmd, error) {
	args := make([]string, 0)
	args = append(args, pc.Args...)
	for _, item := range pc.Flags {
		value, err := envsubst.EvalEnv(item.Value)
		if err != nil {
			return nil, errors.Wrapf(err, `failed to eval flag value "%s" for plugin "%s"`, item.Name, executable)
		}
		args = append(args, fmt.Sprintf(`--%s`, strings.TrimPrefix(item.Name, `--`)))
		args = append(args, value)
	}
	env := make([]string, 0)
	for _, item := range pc.Env {
		value, err := envsubst.EvalEnv(item.Value)
		if err != nil {
			return nil, errors.Wrapf(err, `failed to eval env value "%s" for plugin "%s"`, item.Name, executable)
		}
		env = append(env, fmt.Sprintf(`%s=%s`, strings.ToUpper(item.Name), value))
	}
	/* #nosec */
	cmd := exec.CommandContext(ctx, fmt.Sprintf(`./%s`, executable), args...)
	cmd.Env = append(cmd.Env, env...)
	cmd.Dir = path
	return cmd, nil
}
