package internal

import (
	"context"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/hashicorp/go-plugin"
	"github.com/insidieux/inizio/internal/builtin/layout"
	"github.com/insidieux/inizio/internal/logger"
	"github.com/insidieux/inizio/internal/plugins"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/rhysd/abspath"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type (
	pluginsPath       string
	pluginsConfigPath string

	registryFailFast bool

	layoutCleanup            bool
	layoutTemplateDockerfile *string
	layoutTemplateMakefile   *string
)

func provideCommandViper(set *pflag.FlagSet) (*viper.Viper, error) {
	v := viper.New()
	v.SetEnvPrefix(`INIZIO_`)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	if err := v.BindPFlags(set); err != nil {
		return nil, errors.Wrap(err, `failed to bind commandline flags`)
	}
	return v, nil
}

func provideLoggerLevel(v *viper.Viper) (logrus.Level, error) {
	level := v.GetString(`logger.level`)
	if level == `` {
		level = logrus.InfoLevel.String()
	}
	return logrus.ParseLevel(level)
}

func provideLogger(level logrus.Level) logrus.FieldLogger {
	log := logger.GetLogger()
	log.(*logrus.Logger).SetLevel(level)
	return log
}

func provideTemplate() layout.TemplateInterface {
	return layout.NewTemplate(``, sprig.TxtFuncMap())
}

func provideLayoutCleanup(v *viper.Viper) layoutCleanup {
	cleanup := v.GetBool(`layout.cleanup`)
	return layoutCleanup(cleanup)
}

func provideLayoutTemplateDockerfile(v *viper.Viper) (layoutTemplateDockerfile, error) {
	path := v.GetString(`layout.template.dockerfile`)
	if path == `` {
		return nil, nil
	}
	if !strings.HasSuffix(path, layout.Extension) {
		return nil, errors.Errorf(`layout template dockerfile must end with "%s" extension`, layout.Extension)
	}
	return &path, nil
}

func provideLayoutTemplateMakefile(v *viper.Viper) (layoutTemplateMakefile, error) {
	path := v.GetString(`layout.template.makefile`)
	if path == `` {
		return nil, nil
	}
	if !strings.HasSuffix(path, layout.Extension) {
		return nil, errors.Errorf(`layout template makefile must end with "%s" extension`, layout.Extension)
	}
	return &path, nil
}

func provideGeneratorOptions(cleanup layoutCleanup, dockerfile layoutTemplateDockerfile, makefile layoutTemplateMakefile) []layout.Option {
	return []layout.Option{
		layout.WithCleanup(bool(cleanup)),
		layout.WithTemplateDockerfile(dockerfile),
		layout.WithTemplateMakefile(makefile),
	}
}

func providePluginsConfigPath(v *viper.Viper) pluginsConfigPath {
	return pluginsConfigPath(v.GetString(`plugins.config`))
}

func providePluginsConfig(filesystem afero.Fs, path pluginsConfigPath) (plugins.ConfigInterface, error) {
	config := &plugins.Config{}
	if path != `` {
		contents, err := afero.ReadFile(filesystem, string(path))
		if err != nil {
			return nil, errors.Wrap(err, `failed to read config file`)
		}
		if err := yaml.Unmarshal(contents, config); err != nil {
			return nil, errors.Wrap(err, `failed to parse plugins config file`)
		}
	}
	return config, nil
}

func providePluginsPath(v *viper.Viper) (pluginsPath, error) {
	path := v.GetString(`plugins.path`)
	if path == `` {
		return ``, errors.New(`flag "plugins.path" is required and must be not empty`)
	}
	return pluginsPath(v.GetString(`plugins.path`)), nil
}

func provideRegistryClients(ctx context.Context, loader *plugins.Loader, path pluginsPath) ([]plugins.ClientInterface, error) {
	if path == `` {
		return nil, errors.New(`plugins path parameter is empty`)
	}
	return loader.Load(ctx, string(path))
}

func provideRegistryFailFast(v *viper.Viper) registryFailFast {
	return registryFailFast(v.GetBool(`plugins.fail-fast`))
}

func provideRegistry(clients []plugins.ClientInterface, logger logrus.FieldLogger, failFast registryFailFast) (plugins.RegistryInterface, func()) {
	return plugins.NewRegistry(clients, logger, bool(failFast)), plugin.CleanupClients
}

func provideCore(layout generator.Generator, registry plugins.RegistryInterface, logger logrus.FieldLogger) *core {
	return &core{
		layout:   layout,
		registry: registry,
		logger:   logger,
	}
}

// --- Survey ---

func provideWorkingDirectory(args []string) (string, error) {
	if len(args) == 0 {
		return ``, errors.New(`working directory argument is required`)
	}
	if args[0] == `` {
		return ``, errors.New(`working directory argument cannot be empty`)
	}
	directory, err := abspath.ExpandFrom(args[0])
	if err != nil {
		return ``, errors.Wrapf(err, `failed to get absolute path of working directory`)
	}
	return directory.String(), nil
}
