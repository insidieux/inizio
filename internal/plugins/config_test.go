package plugins

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestConfig(t *testing.T) {
	suite.Run(t, new(configTestSuite))
}

func TestPluginConfig(t *testing.T) {
	suite.Run(t, new(pluginConfigTestSuite))
}

// --- Suites ---

type configTestSuite struct {
	suite.Suite
}

func (s *configTestSuite) TestConfigLookupError() {
	cfg := new(Config)

	lookup, err := cfg.Lookup(`name`)
	s.Error(err)
	s.EqualError(err, `plugin "name" is not defined in config`)
	s.Nil(lookup)
}

func (s *configTestSuite) TestConfigLookup() {
	item := &PluginConfig{Name: `second`}
	cfg := &Config{{Name: `first`}, item}

	lookup, err := cfg.Lookup(`second`)
	s.NoError(err)
	s.NotNil(lookup)
	s.Equal(item, lookup)
}

type pluginConfigTestSuite struct {
	suite.Suite
}

func (s *pluginConfigTestSuite) TestCommandFlagsError() {
	cfg := &PluginConfig{
		Name: `name`,
		Flags: []struct {
			Name  string `json:"name" yaml:"name"`
			Value string `json:"value" yaml:"value"`
		}{
			{
				Name:  `variable`,
				Value: `${VA^R^}`,
			},
		},
	}

	cmd, err := cfg.Command(context.Background(), `/some/path`, `plugin`)
	s.Error(err)
	s.EqualError(err, `failed to eval flag value "variable" for plugin "plugin": bad substitution`)
	s.Nil(cmd)
}

func (s *pluginConfigTestSuite) TestCommandEnvError() {
	cfg := &PluginConfig{
		Name: `name`,
		Env: []struct {
			Name  string `json:"name" yaml:"name"`
			Value string `json:"value" yaml:"value"`
		}{
			{
				Name:  `variable`,
				Value: `${VA^R^}`,
			},
		},
	}

	cmd, err := cfg.Command(context.Background(), `/some/path`, `plugin`)
	s.Error(err)
	s.EqualError(err, `failed to eval env value "variable" for plugin "plugin": bad substitution`)
	s.Nil(cmd)
}

func (s *pluginConfigTestSuite) TestCommand() {
	cfg := &PluginConfig{
		Name: `name`,
		Args: []string{`arg`},
		Env: []struct {
			Name  string `json:"name" yaml:"name"`
			Value string `json:"value" yaml:"value"`
		}{
			{
				Name:  `variable`,
				Value: `${VAR}`,
			},
		},
		Flags: []struct {
			Name  string `json:"name" yaml:"name"`
			Value string `json:"value" yaml:"value"`
		}{
			{
				Name:  `variable`,
				Value: `${VAR}`,
			},
		},
	}

	_ = os.Setenv(`VAR`, `value`)

	cmd, err := cfg.Command(context.Background(), `/some/path`, `plugin`)
	s.NoError(err)
	s.NotNil(cmd)
	s.Equal(`/some/path`, cmd.Dir)
	s.Equal(`./plugin`, cmd.Path)
	s.Equal([]string{`VARIABLE=value`}, cmd.Env)
	s.Equal([]string{`./plugin`, `arg`, `--variable`, `value`}, cmd.Args)
}
