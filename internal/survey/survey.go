package survey

import (
	"path/filepath"
	"runtime"
	"strings"

	surveySDK "github.com/AlecAivazis/survey/v2"

	"github.com/blang/semver/v4"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"golang.org/x/mod/module"
)

type (
	answers struct {
		ApplicationName        string `survey:"application.name"`
		ApplicationDescription string `survey:"application.description"`
		GoModule               string `survey:"golang.module"`
		GoVersion              string `survey:"golang.version"`
		Override               bool   `survey:"override"`
		workingDirectory       string `survey:"-"`
	}
)

// MakeRunOptions cast itself to generator.RunOptions.
func (a *answers) MakeRunOptions() generator.RunOptions {
	return generator.RunOptions{
		Override:         a.Override,
		WorkingDirectory: a.workingDirectory,
	}
}

// MakeRunValues cast itself to generator.RunValues, which will be used in generation process.
func (a *answers) MakeRunValues() generator.RunValues {
	return generator.RunValues{
		Application: generator.RunValuesApplication{
			Name:        a.ApplicationName,
			Description: a.ApplicationDescription,
		},
		Golang: generator.RunValuesGolang{
			Module:  a.GoModule,
			Version: a.GoVersion,
		},
	}
}

// Ask start tty terminal with series of prompt for collection data for generation of project.
func Ask(directory string) (generator.RunOptions, generator.RunValues, error) {
	base := filepath.Base(directory)
	version := strings.TrimPrefix(runtime.Version(), `go`)
	questions := []*surveySDK.Question{
		{
			Name: `application.name`,
			Prompt: &surveySDK.Input{
				Message: `Provide application name`,
				Default: base,
			},
			Validate: surveySDK.Required,
			Transform: func(ans interface{}) (newAns interface{}) {
				return strings.TrimSpace(cast.ToString(ans))
			},
		},
		{
			Name: `application.description`,
			Prompt: &surveySDK.Input{
				Message: `Provide application description`,
			},
			Transform: func(ans interface{}) (newAns interface{}) {
				return strings.TrimSpace(cast.ToString(ans))
			},
		},
		{
			Name: `golang.module`,
			Prompt: &surveySDK.Input{
				Message: `Provide golang module name`,
			},
			Validate: surveySDK.ComposeValidators(
				surveySDK.Required,
				func(answer interface{}) error {
					if err := module.CheckPath(cast.ToString(answer)); err != nil {
						return errors.Wrap(err, `failed to check module path`)
					}
					return nil
				},
			),
			Transform: func(ans interface{}) (newAns interface{}) {
				return strings.TrimSpace(cast.ToString(ans))
			},
		},
		{
			Name: `golang.version`,
			Prompt: &surveySDK.Input{
				Message: `Provide golang version`,
				Default: version,
			},
			Validate: surveySDK.ComposeValidators(
				surveySDK.Required,
				func(answer interface{}) error {
					version, err := semver.ParseTolerant(cast.ToString(answer))
					if err != nil {
						return errors.Wrap(err, `failed to parse golang version`)
					}
					if err := version.Validate(); err != nil {
						return errors.Wrap(err, `golang version validation failed`)
					}
					return nil
				},
			),
			Transform: func(ans interface{}) (newAns interface{}) {
				return strings.TrimSpace(cast.ToString(ans))
			},
		},
		{
			Name: `override`,
			Prompt: &surveySDK.Confirm{
				Message: `Overwrite existing files?`,
				Default: true,
			},
			Validate: surveySDK.Required,
		},
	}
	response := new(answers)
	response.workingDirectory = directory
	if err := surveySDK.Ask(questions, response); err != nil {
		return generator.RunOptions{}, generator.RunValues{}, errors.Wrap(err, `failed to ask application name`)
	}
	return response.MakeRunOptions(), response.MakeRunValues(), nil
}
