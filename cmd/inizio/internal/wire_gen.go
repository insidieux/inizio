// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package internal

import (
	"context"
	"github.com/insidieux/inizio/internal/builtin/layout"
	"github.com/insidieux/inizio/internal/plugins"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

// Injectors from wire.go:

func newCore(contextContext context.Context, flagSet *pflag.FlagSet) (*core, func(), error) {
	boxDirectory := _wireBoxDirectoryValue
	readFileFS, err := layout.ProvideEmbedFS(boxDirectory)
	if err != nil {
		return nil, nil, err
	}
	fs := afero.NewOsFs()
	boxInterface := layout.NewBox(readFileFS, fs)
	templateInterface := provideTemplate()
	rendererInterface := layout.NewRenderer(boxInterface, templateInterface)
	viper, err := provideCommandViper(flagSet)
	if err != nil {
		return nil, nil, err
	}
	internalLayoutCleanup := provideLayoutCleanup(viper)
	internalLayoutTemplateDockerfile, err := provideLayoutTemplateDockerfile(viper)
	if err != nil {
		return nil, nil, err
	}
	internalLayoutTemplateMakefile, err := provideLayoutTemplateMakefile(viper)
	if err != nil {
		return nil, nil, err
	}
	v := provideGeneratorOptions(internalLayoutCleanup, internalLayoutTemplateDockerfile, internalLayoutTemplateMakefile)
	generator := layout.NewGenerator(rendererInterface, fs, v...)
	internalPluginsConfigPath := providePluginsConfigPath(viper)
	configInterface, err := providePluginsConfig(fs, internalPluginsConfigPath)
	if err != nil {
		return nil, nil, err
	}
	level, err := provideLoggerLevel(viper)
	if err != nil {
		return nil, nil, err
	}
	fieldLogger := provideLogger(level)
	loader := plugins.NewLoader(configInterface, fs, fieldLogger)
	internalPluginsPath, err := providePluginsPath(viper)
	if err != nil {
		return nil, nil, err
	}
	v2, err := provideRegistryClients(contextContext, loader, internalPluginsPath)
	if err != nil {
		return nil, nil, err
	}
	internalRegistryFailFast := provideRegistryFailFast(viper)
	registryInterface, cleanup := provideRegistry(v2, fieldLogger, internalRegistryFailFast)
	internalCore := provideCore(generator, registryInterface, fieldLogger)
	return internalCore, func() {
		cleanup()
	}, nil
}

var (
	_wireBoxDirectoryValue = layout.BoxDirectory(layout.EmbedDirectory)
)
