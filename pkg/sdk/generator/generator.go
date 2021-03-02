package generator

import (
	"context"
)

type (
	// RunOptions contains common options for generation process
	RunOptions struct {
		Override         bool
		WorkingDirectory string
	}

	// RunValues contains information for newly generated project such as RunValuesApplication and RunValuesGolang
	RunValues struct {
		Application RunValuesApplication
		Golang      RunValuesGolang
	}

	// RunValuesApplication contains information about name and description of project.
	RunValuesApplication struct {
		Name        string
		Description string
	}

	// RunValuesGolang contains information about go module name and used golang version.
	RunValuesGolang struct {
		Module  string
		Version string
	}

	// RunResult contains the slice of string representing list of generated directories and files.
	RunResult []string

	// Generator is the interface that we are exposing as main for writing plugins.
	Generator interface {
		Run(context.Context, RunOptions, RunValues) (RunResult, error)
	}
)
