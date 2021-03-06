package layout

import (
	"embed"
	"io/fs"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const (
	// EmbedDirectory is main embed directory with predefined templates.
	EmbedDirectory = `embed`

	// Extension is a common extension for all template files used for generation project.
	Extension = `gotmpl`

	fileNameGoMod        = "go.mod"
	fileNameMainGo       = "main.go"
	fileNameReadmeMD     = "README.md"
	fileNameDockerfile   = "Dockerfile"
	fileNameMakefile     = "Makefile"
	fileNameDockerIgnore = ".dockerignore"
)

type (
	// BoxDirectory is an alias for string type, determine where placed builtin layout templates.
	BoxDirectory string

	// BoxInterface is a common interface for read file content of templates, builtin and passed by values/args.
	BoxInterface interface {
		ReadFile(string) ([]byte, error)
	}

	// Box is a BoxInterface implementation. Contains information about embed and os filesystem.
	Box struct {
		embedFS fs.ReadFileFS
		osFS    afero.Fs
	}
)

var (
	//go:embed embed/*
	filesystem embed.FS

	_ BoxInterface = new(Box)
)

// ProvideEmbedFS return fs.ReadFileFS, which contains embed built-in template files.
func ProvideEmbedFS(directory BoxDirectory) (fs.ReadFileFS, error) {
	sub, err := fs.Sub(filesystem, string(directory))
	if err != nil {
		return nil, errors.Wrapf(err, `failed to open directory "%s"`, directory)
	}
	return sub.(fs.ReadFileFS), nil
}

// NewBox return BoxInterface, which contains embed built-in template files.
func NewBox(embedFS fs.ReadFileFS, osFS afero.Fs) BoxInterface {
	return &Box{
		embedFS: embedFS,
		osFS:    osFS,
	}
}

// ReadFile implements BoxInterface method. Read from embed FS, then os FS.
func (b *Box) ReadFile(name string) ([]byte, error) {
	content, err := b.embedFS.ReadFile(name)
	if err == nil {
		return content, nil
	}
	content, err = afero.ReadFile(b.osFS, name)
	if err == nil {
		return content, nil
	}
	return nil, errors.Wrapf(err, `failed to read file "%s"`, name)
}
