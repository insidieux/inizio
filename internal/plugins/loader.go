package plugins

import (
	"context"
	"regexp"

	internalSDK "github.com/insidieux/inizio/pkg/sdk/generator/plugin"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type (
	// Loader prepare slice of ClientInterface used by RegistryInterface.
	Loader struct {
		config     ConfigInterface
		filesystem afero.Fs
		logger     logrus.FieldLogger
	}
)

const (
	pluginBinaryRegex = "^inizio-plugin-(.*)$"
)

var (
	nameRegexp = regexp.MustCompile(pluginBinaryRegex)
)

// NewLoader create new loader with readonly and regexp filesystem for searching suitable plugins.
func NewLoader(config ConfigInterface, filesystem afero.Fs, logger logrus.FieldLogger) *Loader {
	return &Loader{
		config:     config,
		filesystem: afero.NewRegexpFs(afero.NewReadOnlyFs(filesystem), nameRegexp),
		logger:     logger,
	}
}

// Load use passed ConfigInterface and afero.Fs for find all suitable plugin binaries in called path.
func (l *Loader) Load(ctx context.Context, path string) ([]ClientInterface, error) {
	list, err := afero.ReadDir(l.filesystem, path)
	if err != nil {
		return nil, errors.Wrapf(err, `failed to read plugin directory "%s"`, path)
	}

	clients := make([]ClientInterface, 0)
	for _, item := range list {
		parts := nameRegexp.FindStringSubmatch(item.Name())
		if len(parts) < 2 {
			continue
		}
		name := parts[1]
		cfg, err := l.config.Lookup(name)
		if err != nil {
			cfg = newPluginConfig(name)
		}
		executable, err := cfg.Command(ctx, path, item.Name())
		if err != nil {
			return nil, errors.Wrapf(err, `failed to make os/exec command for plugin "%s"`, item.Name())
		}
		clients = append(
			clients,
			&client{
				name:       item.Name(),
				downstream: internalSDK.NewClient(executable),
			},
		)
	}
	return clients, nil
}
