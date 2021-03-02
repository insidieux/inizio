package layout

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestLayoutOptions(t *testing.T) {
	suite.Run(t, new(layoutOptionsTestSuite))
}

func TestLayoutGenerator(t *testing.T) {
	suite.Run(t, new(layoutGeneratorTestSuite))
}

// --- Suites ---

type layoutOptionsTestSuite struct {
	suite.Suite
}

func (s *layoutOptionsTestSuite) TestOptionWithCleanup() {
	options := &Options{}

	WithCleanup(true)(options)
	s.True(options.Cleanup)
}

func (s *layoutOptionsTestSuite) TestOptionWithTemplateDockerfile() {
	options := &Options{}
	path := `/some/path`

	WithTemplateDockerfile(&path)(options)
	s.Equal(path, *options.Templates.Dockerfile)
}

func (s *layoutOptionsTestSuite) TestOptionWithTemplateMakefile() {
	options := &Options{}
	path := `/some/path`

	WithTemplateMakefile(&path)(options)
	s.Equal(path, *options.Templates.Makefile)
}

type layoutGeneratorTestSuite struct {
	suite.Suite
	target  file
	options generator.RunOptions
	values  generator.RunValues
}

func (s *layoutGeneratorTestSuite) SetupTest() {
	s.target = file{
		readPath:  `/read/path/file.gotmpl`,
		writePath: `/write/path/file.go`,
	}
	s.options = generator.RunOptions{
		Override:         false,
		WorkingDirectory: `/working/directory`,
	}
	s.values = generator.RunValues{
		Application: generator.RunValuesApplication{
			Name:        `application`,
			Description: `Short description`,
		},
		Golang: generator.RunValuesGolang{
			Module:  `github.com/insidieux/example`,
			Version: `1.15`,
		},
	}
}

func (s *layoutGeneratorTestSuite) TearDownTest() {
	s.options.Override = false
}

func (s *layoutGeneratorTestSuite) TestProcessFileRenderError() {
	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.
		On(`Render`, s.target.readPath, s.values).
		Return(nil, errors.New(`expected error`))

	gen := NewGenerator(mockedRenderer, nil)
	err := gen.processFile(s.target, s.options, s.values)
	s.Error(err)
	s.EqualError(err, `failed to render template: expected error`)
}

func (s *layoutGeneratorTestSuite) TestProcessFileMkdirError() {
	writeDirectory := filepath.Dir(s.target.writePath)

	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.On(`Render`, s.target.readPath, s.values).Return([]byte(`content`), nil)

	mockedFilesystem := new(mockFilesystem)
	mockedFilesystem.On(`MkdirAll`, writeDirectory, mock.Anything).Return(errors.New(`expected error`))

	gen := NewGenerator(mockedRenderer, mockedFilesystem)
	err := gen.processFile(s.target, s.options, s.values)
	s.Error(err)
	s.EqualError(err, `failed to create directory "/write/path": expected error`)
}

func (s *layoutGeneratorTestSuite) TestProcessFileExistedFileError() {
	writeDirectory := filepath.Dir(s.target.writePath)

	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.On(`Render`, s.target.readPath, s.values).Return([]byte(`content`), nil)

	filesystem := afero.NewMemMapFs()
	_ = filesystem.MkdirAll(writeDirectory, 0755)
	_, _ = filesystem.Create(s.target.writePath)

	gen := NewGenerator(mockedRenderer, filesystem)
	err := gen.processFile(s.target, s.options, s.values)
	s.Error(err)
	s.EqualError(err, `file "/write/path/file.go" is already exists, but override option is not passed`)
}

func (s *layoutGeneratorTestSuite) TestProcessFileExistedRemoveError() {
	writeDirectory := filepath.Dir(s.target.writePath)

	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.On(`Render`, s.target.readPath, s.values).Return([]byte(`content`), nil)

	filesystem := afero.NewMemMapFs()
	_ = filesystem.MkdirAll(writeDirectory, 0755)
	memFile, _ := filesystem.Create(s.target.writePath)
	memFileInfo, _ := memFile.Stat()

	mockedFilesystem := new(mockFilesystem)
	mockedFilesystem.On(`MkdirAll`, writeDirectory, mock.Anything).Return(nil)
	mockedFilesystem.On(`Stat`, s.target.writePath).Return(memFileInfo, nil)
	mockedFilesystem.On(`Remove`, s.target.writePath).Return(errors.New(`expected error`))

	s.options.Override = true

	gen := NewGenerator(mockedRenderer, mockedFilesystem)
	err := gen.processFile(s.target, s.options, s.values)
	s.Error(err)
	s.EqualError(err, `failed to remove previous file "/write/path/file.go": expected error`)
}

func (s *layoutGeneratorTestSuite) TestProcessFileCreateError() {
	writeDirectory := filepath.Dir(s.target.writePath)

	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.On(`Render`, s.target.readPath, s.values).Return([]byte(`content`), nil)

	mockedFilesystem := new(mockFilesystem)
	mockedFilesystem.On(`MkdirAll`, writeDirectory, mock.Anything).Return(nil)
	mockedFilesystem.On(`Stat`, s.target.writePath).Return(nil, nil)
	mockedFilesystem.On(`Create`, s.target.writePath).Return(nil, errors.New(`expected error`))

	gen := NewGenerator(mockedRenderer, mockedFilesystem)
	err := gen.processFile(s.target, s.options, s.values)
	s.Error(err)
	s.EqualError(err, `failed to create file "/write/path/file.go": expected error`)
}

func (s *layoutGeneratorTestSuite) TestProcessFileWriteError() {
	content := []byte(`content`)
	writeDirectory := filepath.Dir(s.target.writePath)

	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.On(`Render`, s.target.readPath, s.values).Return(content, nil)

	mockedFile := new(mockFile)
	mockedFile.On(`Write`, content).Return(0, errors.New(`expected error`))

	mockedFilesystem := new(mockFilesystem)
	mockedFilesystem.On(`MkdirAll`, writeDirectory, mock.Anything).Return(nil)
	mockedFilesystem.On(`Stat`, s.target.writePath).Return(nil, nil)
	mockedFilesystem.On(`Create`, s.target.writePath).Return(mockedFile, nil)

	gen := NewGenerator(mockedRenderer, mockedFilesystem)
	err := gen.processFile(s.target, s.options, s.values)
	s.Error(err)
	s.EqualError(err, `failed to write to file "/write/path/file.go": expected error`)
}

func (s *layoutGeneratorTestSuite) TestProcessFile() {
	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.On(`Render`, s.target.readPath, s.values).Return([]byte(`content`), nil)

	filesystem := afero.NewMemMapFs()

	gen := NewGenerator(mockedRenderer, filesystem)
	err := gen.processFile(s.target, s.options, s.values)
	s.NoError(err)
}

func (s *layoutGeneratorTestSuite) TestPrepareFilesList() {
	gen := NewGenerator(nil, nil)
	items := gen.prepareFilesList(s.options, s.values)
	s.NotEmpty(items)
	s.Len(items, 6)
	s.Equal(filepath.Join(s.options.WorkingDirectory, fileNameGoMod), items[0].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, subPathMainGo, strings.ToLower(s.values.Application.Name), fileNameMainGo), items[1].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, fileNameReadmeMD), items[2].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, fileNameDockerIgnore), items[3].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, fileNameMakefile), items[4].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, subPathDockerfile, strings.ToLower(s.values.Application.Name), fileNameDockerfile), items[5].writePath)
	s.Equal(fmt.Sprintf(`%s.%s`, fileNameMakefile, Extension), items[4].readPath)
	s.Equal(fmt.Sprintf(`%s.%s`, fileNameDockerfile, Extension), items[5].readPath)
}

func (s *layoutGeneratorTestSuite) TestPrepareFilesListWithRedefine() {
	dockerfilePath := `/template/path/Dockerfile.gotmpl`
	makefilePath := `/template/path/Dockerfile.gotmpl`

	gen := NewGenerator(nil, nil, WithTemplateDockerfile(&dockerfilePath), WithTemplateMakefile(&makefilePath))
	items := gen.prepareFilesList(s.options, s.values)
	s.NotEmpty(items)
	s.Len(items, 6)
	s.Equal(filepath.Join(s.options.WorkingDirectory, fileNameGoMod), items[0].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, subPathMainGo, strings.ToLower(s.values.Application.Name), fileNameMainGo), items[1].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, fileNameReadmeMD), items[2].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, fileNameDockerIgnore), items[3].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, fileNameMakefile), items[4].writePath)
	s.Equal(filepath.Join(s.options.WorkingDirectory, subPathDockerfile, strings.ToLower(s.values.Application.Name), fileNameDockerfile), items[5].writePath)
	s.Equal(dockerfilePath, items[4].readPath)
	s.Equal(makefilePath, items[5].readPath)
}

func (s *layoutGeneratorTestSuite) TestRunCleanupError() {
	mockedFilesystem := new(mockFilesystem)
	mockedFilesystem.On(`RemoveAll`, s.options.WorkingDirectory).Return(errors.New(`expected error`))

	gen := NewGenerator(nil, mockedFilesystem, WithCleanup(true))
	result, err := gen.Run(context.Background(), s.options, s.values)
	s.Error(err)
	s.EqualError(err, `failed to cleanup "/working/directory" directory: expected error`)
	s.Nil(result)
}

func (s *layoutGeneratorTestSuite) TestRunProcessFileError() {
	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.
		On(`Render`, fmt.Sprintf(`%s.%s`, fileNameGoMod, Extension), s.values).
		Return(nil, errors.New(`expected error`))

	filesystem := afero.NewMemMapFs()

	gen := NewGenerator(mockedRenderer, filesystem, WithCleanup(true))
	result, err := gen.Run(context.Background(), s.options, s.values)
	s.Error(err)
	s.EqualError(err, `failed to process template "go.mod.gotmpl": failed to render template: expected error`)
	s.Nil(result)
}

func (s *layoutGeneratorTestSuite) TestRun() {
	mockedRenderer := new(mockRendererInterface)
	mockedRenderer.
		On(`Render`, mock.Anything, s.values).
		Return([]byte(`content`), nil)

	filesystem := afero.NewMemMapFs()

	gen := NewGenerator(mockedRenderer, filesystem, WithCleanup(true))
	result, err := gen.Run(context.Background(), s.options, s.values)
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 6)
	s.Subset(
		result,
		[]string{
			`/working/directory/go.mod`,
			`/working/directory/cmd/application/main.go`,
			`/working/directory/README.md`,
			`/working/directory/.dockerignore`,
			`/working/directory/Makefile`,
			`/working/directory/build/docker/cmd/application/Dockerfile`,
		},
	)
}
