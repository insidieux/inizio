package layout

import (
	"io/fs"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestProvideEmbedFS(t *testing.T) {
	suite.Run(t, new(provideEmbedFSTestSuite))
}

func TestBox(t *testing.T) {
	suite.Run(t, new(boxTestSuite))
}

// --- Suites ---

type provideEmbedFSTestSuite struct {
	suite.Suite
}

func (s *provideEmbedFSTestSuite) TestSubError() {
	filesystem, err := ProvideEmbedFS(`..`)
	s.Error(err)
	s.Nil(filesystem)
}

func (s *provideEmbedFSTestSuite) TestSuccess() {
	filesystem, err := ProvideEmbedFS(EmbedDirectory)
	s.NoError(err)
	s.Implements(new(fs.ReadFileFS), filesystem)
}

type boxTestSuite struct {
	suite.Suite

	embedFS fs.ReadFileFS
	osFS    afero.Fs
}

func (s *boxTestSuite) SetupTest() {
	s.embedFS, _ = ProvideEmbedFS(EmbedDirectory)
	s.osFS = afero.NewMemMapFs()
}

func (s *boxTestSuite) TestEmbedReadFile() {
	box := NewBox(s.embedFS, s.osFS)
	expected := []byte(`package main

func main() {

}
`)

	content, err := box.ReadFile(`main.go.gotmpl`)
	s.NoError(err)
	s.Equal(expected, content)
}

func (s *boxTestSuite) TestOSReadFile() {
	expected := []byte(`package main

func main() {

}
`)
	f, _ := s.osFS.Create(`/some/path/main.go.gotmpl`)
	_, _ = f.Write(expected)
	box := NewBox(s.embedFS, s.osFS)

	content, err := box.ReadFile(`/some/path/main.go.gotmpl`)
	s.NoError(err)
	s.Equal(expected, content)
}

func (s *boxTestSuite) TestError() {
	box := NewBox(s.embedFS, s.osFS)

	content, err := box.ReadFile(`/some/path/main.go.gotmpl`)
	s.Error(err)
	s.EqualError(err, `failed to read file "/some/path/main.go.gotmpl": open /some/path/main.go.gotmpl: file does not exist`)
	s.Nil(content)
}
