package internal

import (
	"context"

	"github.com/insidieux/inizio/internal/plugins"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	core struct {
		layout   generator.Generator
		registry plugins.RegistryInterface
		logger   logrus.FieldLogger
	}
)

func (c *core) Run(ctx context.Context, options generator.RunOptions, values generator.RunValues) error {
	c.logger.Infoln(`Start layout generation`)
	result, err := c.layout.Run(ctx, options, values)
	if err != nil {
		return errors.Wrap(err, `failed to generate layout`)
	}
	for _, item := range result {
		c.logger.Infof(`Layout generated: %s`, item)
	}
	if c.registry.HasPlugins() {
		c.logger.Infoln(`Start plugins generation`)
		result, err = c.registry.Process(ctx, options, values)
		if err != nil {
			return errors.Wrap(err, `failed to process process plugins`)
		}
		for _, item := range result {
			c.logger.Infof(`Plugin generated: %s`, item)
		}
	}
	return nil
}
