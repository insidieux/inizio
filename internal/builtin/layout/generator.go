package layout

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const (
	subPathMainGo     = `cmd`
	subPathDockerfile = `build/docker/cmd`
)

type (
	// Generator is built-in generator.Generator implementation. Used to generate common/base layout for project.
	Generator struct {
		renderer   RendererInterface
		filesystem afero.Fs
		options    *Options
	}

	// Options provides additional Generator settings, such as generator Cleanup or Templates.
	Options struct {
		Cleanup   bool
		Templates struct {
			Dockerfile *string
			Makefile   *string
		}
	}

	// Option allows setting optional Generator Options.
	Option func(options *Options)

	file struct {
		writePath string
		readPath  string
	}
)

var (
	_ generator.Generator = &Generator{}
)

// WithCleanup specifies Options.Cleanup value, allow cleanup working directory before generation start.
func WithCleanup(cleanup bool) Option {
	return func(options *Options) {
		options.Cleanup = cleanup
	}
}

// WithTemplateDockerfile specifies custom Dockerfile template path.
func WithTemplateDockerfile(path *string) Option {
	return func(options *Options) {
		options.Templates.Dockerfile = path
	}
}

// WithTemplateMakefile specifies custom Makefile template path.
func WithTemplateMakefile(path *string) Option {
	return func(options *Options) {
		options.Templates.Makefile = path
	}
}

// NewGenerator return Generator
func NewGenerator(renderer RendererInterface, filesystem afero.Fs, options ...Option) *Generator {
	gen := &Generator{
		renderer:   renderer,
		filesystem: filesystem,
		options:    new(Options),
	}
	for _, option := range options {
		option(gen.options)
	}
	return gen
}

/*
Run performs main process of generation.
Current steps:
 - make cleanup working directory if called
 - prepare file list for generation
 - process each file/template
*/
func (g *Generator) Run(_ context.Context, options generator.RunOptions, values generator.RunValues) (generator.RunResult, error) {
	if g.options.Cleanup {
		if err := g.filesystem.RemoveAll(options.WorkingDirectory); err != nil {
			return nil, errors.Wrapf(err, `failed to cleanup "%s" directory`, options.WorkingDirectory)
		}
	}
	result := make([]string, 0)
	for _, f := range g.prepareFilesList(options, values) {
		if err := g.processFile(f, options, values); err != nil {
			return nil, errors.Wrapf(err, `failed to process template "%s"`, f.readPath)
		}
		result = append(result, f.writePath)
	}
	return result, nil
}

func (g *Generator) prepareFilesList(options generator.RunOptions, values generator.RunValues) []file {
	var files = []file{
		{
			writePath: filepath.Join(options.WorkingDirectory, fileNameGoMod),
			readPath:  fmt.Sprintf(`%s.%s`, fileNameGoMod, Extension),
		},
		{
			writePath: filepath.Join(options.WorkingDirectory, subPathMainGo, strings.ToLower(values.Application.Name), fileNameMainGo),
			readPath:  fmt.Sprintf(`%s.%s`, fileNameMainGo, Extension),
		},
		{
			writePath: filepath.Join(options.WorkingDirectory, fileNameReadmeMD),
			readPath:  fmt.Sprintf(`%s.%s`, fileNameReadmeMD, Extension),
		},
		{
			writePath: filepath.Join(options.WorkingDirectory, fileNameDockerIgnore),
			readPath:  fmt.Sprintf(`%s.%s`, fileNameDockerIgnore, Extension),
		},
	}
	// Makefile override
	makefileTemplate := fmt.Sprintf(`%s.%s`, fileNameMakefile, Extension)
	if g.options.Templates.Makefile != nil {
		makefileTemplate = *g.options.Templates.Makefile
	}
	files = append(files, file{
		writePath: filepath.Join(options.WorkingDirectory, fileNameMakefile),
		readPath:  makefileTemplate,
	})
	// Dockerfile override
	dockerfileTemplate := fmt.Sprintf(`%s.%s`, fileNameDockerfile, Extension)
	if g.options.Templates.Dockerfile != nil {
		dockerfileTemplate = *g.options.Templates.Dockerfile
	}
	files = append(files, file{
		writePath: filepath.Join(options.WorkingDirectory, subPathDockerfile, strings.ToLower(values.Application.Name), fileNameDockerfile),
		readPath:  dockerfileTemplate,
	})
	return files
}

func (g *Generator) processFile(f file, options generator.RunOptions, values generator.RunValues) error {
	content, err := g.renderer.Render(f.readPath, values)
	if err != nil {
		return errors.Wrap(err, `failed to render template`)
	}
	writeDirectory := filepath.Dir(f.writePath)
	if err = g.filesystem.MkdirAll(writeDirectory, 0755); err != nil {
		return errors.Wrapf(err, `failed to create directory "%s"`, writeDirectory)
	}
	fi, _ := g.filesystem.Stat(f.writePath)
	if fi != nil {
		if !options.Override {
			return errors.Errorf(`file "%s" is already exists, but override option is not passed`, f.writePath)
		}
		if err := g.filesystem.Remove(f.writePath); err != nil {
			return errors.Wrapf(err, `failed to remove previous file "%s"`, f.writePath)
		}
	}
	d, err := g.filesystem.Create(f.writePath)
	if err != nil {
		return errors.Wrapf(err, `failed to create file "%s"`, f.writePath)
	}
	if _, err := d.Write(content); err != nil {
		return errors.Wrapf(err, `failed to write to file "%s"`, f.writePath)
	}
	return nil
}
