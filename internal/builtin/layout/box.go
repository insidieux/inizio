package layout

import (
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
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

// NewBox return packd.Box, which contains embed built-in template files
func NewBox() packd.Box {
	return packr.New("templates", "./embed")
}
