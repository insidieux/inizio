package layout

import (
	"embed"
	"io/fs"
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

// NewBox return embed.FS, which contains embed built-in template files
func NewBox() fs.ReadFileFS {
	return filesystem
}
