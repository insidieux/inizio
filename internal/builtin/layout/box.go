package layout

import (
	"embed"
	"io/fs"

	"github.com/pkg/errors"
)

const (
	// Extension is a common extension for all template files used for generation project
	Extension = `gotmpl`

	fileNameGoMod        = "go.mod"
	fileNameMainGo       = "main.go"
	fileNameReadmeMD     = "README.md"
	fileNameDockerfile   = "Dockerfile"
	fileNameMakefile     = "Makefile"
	fileNameDockerIgnore = ".dockerignore"
)

//go:embed embed/*
var filesystem embed.FS

// NewBox return fs.ReadFileFS, which contains embed built-in template files
func NewBox() (fs.ReadFileFS, error) {
	sub, err := fs.Sub(filesystem, `embed`)
	if err != nil {
		return nil, errors.Wrap(err, `failed to open embed directory`)
	}
	return sub.(fs.ReadFileFS), err
}
