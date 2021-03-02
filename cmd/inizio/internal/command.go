package internal

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/insidieux/inizio/internal/builtin/layout"
	"github.com/insidieux/inizio/internal/logger"
	"github.com/insidieux/inizio/internal/survey"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand create cobra.Command for main process
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `inizio`,
		Short:         `Golang boilerplate project generator`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return errors.New(`working directory argument is required`)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.GetLogger().Infoln(`Initialize backend...`)
			c, cleanup, err := newCore(cmd.Context(), cmd.Flags())
			if cleanup != nil {
				defer cleanup()
			}
			if err != nil {
				return errors.Wrap(err, `failed to bootstrap backend`)
			}
			logger.GetLogger().Infoln(`Initialize user survey...`)
			directory, err := provideWorkingDirectory(args)
			if err != nil {
				return errors.Wrap(err, `failed to bootstrap survey`)
			}
			options, values, err := survey.Ask(directory)
			if err != nil {
				if !errors.Is(err, terminal.InterruptErr) {
					return errors.Wrap(err, `failed to handle user survey`)
				}
				return nil
			}
			if err = c.Run(cmd.Context(), options, values); err != nil {
				return errors.Wrap(err, `failed to generate project`)
			}
			return nil
		},
	}
	cmd.Flags().String(`plugins.path`, `/usr/local/bin/inizio-plugins`, `path to plugins directory`)
	cmd.Flags().String(`plugins.config`, ``, `path to plugins config yaml file`)
	cmd.Flags().Bool(`plugins.fail-fast`, false, `stop after first plugin failure`)
	cmd.Flags().Bool(`layout.cleanup`, false, `cleanup working directory before generation`)
	cmd.Flags().String(`layout.template.dockerfile`, ``, fmt.Sprintf(`path to custom Dockerfile template (must have "%s" extension)`, layout.Extension))
	cmd.Flags().String(`layout.template.makefile`, ``, fmt.Sprintf(`path to custom Makefile template (must have "%s" extension)`, layout.Extension))
	cmd.Flags().String(`logger.level`, logrus.InfoLevel.String(), `log level`)
	return cmd
}
