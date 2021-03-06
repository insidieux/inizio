// +build wireinject

package internal

import (
	"context"

	"github.com/google/wire"
	"github.com/insidieux/inizio/internal/builtin/layout"
	"github.com/insidieux/inizio/internal/plugins"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

func newCore(context.Context, *pflag.FlagSet) (*core, func(), error) {
	panic(wire.Build(
		provideCommandViper,
		afero.NewOsFs,
		wire.NewSet(
			provideLoggerLevel,
			provideLogger,
		),
		wire.NewSet(
			wire.NewSet(
				layout.NewBox,
				provideTemplate,
				layout.NewRenderer,
			),
			wire.NewSet(
				provideLayoutCleanup,
				provideLayoutTemplateDockerfile,
				provideLayoutTemplateMakefile,
				provideGeneratorOptions,
			),
			wire.NewSet(
				layout.NewGenerator,
				wire.Bind(new(generator.Generator), new(*layout.Generator)),
			),
		),
		wire.NewSet(
			wire.NewSet(
				wire.NewSet(
					wire.NewSet(
						providePluginsConfigPath,
						providePluginsConfig,
					),
					plugins.NewLoader,
				),
				providePluginsPath,
				provideRegistryClients,
			),
			provideRegistryFailFast,
			provideRegistry,
		),
		provideCore,
	))
}
