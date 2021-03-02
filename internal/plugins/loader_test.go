package plugins

import (
	"context"
	"os/exec"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestLoader(t *testing.T) {
	suite.Run(t, new(loaderTestSuite))
}

// --- Suites ---

type loaderTestSuite struct {
	suite.Suite
}

func (s *loaderTestSuite) TestLoadReadDirectoryError() {
	mockedConfig := new(mockConfigInterface)
	filesystem := afero.NewMemMapFs()

	loader := NewLoader(mockedConfig, filesystem, logrus.New())
	clients, err := loader.Load(context.Background(), `/path`)
	s.Error(err)
	s.EqualError(err, `failed to read plugin directory "/path": open /path: file does not exist`)
	s.Nil(clients)
}

func (s *loaderTestSuite) TestLoadCommandError() {
	ctx := context.Background()
	path := `/path`

	mockedPluginConfig := new(mockPluginConfigInterface)
	mockedPluginConfig.
		On(`Command`, ctx, path, `inizio-plugin-first`).
		Return(nil, errors.New(`expected error`))

	mockedConfig := new(mockConfigInterface)
	mockedConfig.On(`Lookup`, `first`).Return(mockedPluginConfig, nil)

	filesystem := afero.NewMemMapFs()
	_ = filesystem.MkdirAll(`/path`, 0755)
	_, _ = filesystem.Create(`/path/inizio-plugin-first`)

	loader := NewLoader(mockedConfig, filesystem, logrus.New())
	clients, err := loader.Load(ctx, path)
	s.Error(err)
	s.EqualError(err, `failed to make os/exec command for plugin "inizio-plugin-first": expected error`)
	s.Nil(clients)
}

func (s *loaderTestSuite) TestLoad() {
	ctx := context.Background()
	path := `/path`

	mockedPluginConfig := new(mockPluginConfigInterface)
	mockedPluginConfig.
		On(`Command`, ctx, path, `inizio-plugin-second`).
		Return(exec.Command(`/path/inizio-plugin-second`), nil)

	mockedConfig := new(mockConfigInterface)
	mockedConfig.On(`Lookup`, `first`).Return(nil, errors.New(`plugin "second" is not defined in config`))
	mockedConfig.On(`Lookup`, `second`).Return(mockedPluginConfig, nil)

	filesystem := afero.NewMemMapFs()
	_ = filesystem.MkdirAll(`/path`, 0755)
	_, _ = filesystem.Create(`/path/inizio-plugin-first`)
	_, _ = filesystem.Create(`/path/inizio-plugin-second`)

	loader := NewLoader(mockedConfig, filesystem, logrus.New())
	clients, err := loader.Load(ctx, path)
	s.NoError(err)
	s.NotNil(clients)
	s.Len(clients, 2)
}
