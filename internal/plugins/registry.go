package plugins

import (
	"context"
	"fmt"
	"unicode"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

type (
	// RegistryInterface interface provide information how many plugins was registered and how to process them.
	RegistryInterface interface {
		HasPlugins() bool
		Process(context.Context, generator.RunOptions, generator.RunValues) (generator.RunResult, error)
	}

	// Registry is RegistryInterface implementation.
	Registry struct {
		clients  []ClientInterface
		logger   logrus.FieldLogger
		failFast bool
	}
)

var (
	_ RegistryInterface = &Registry{}
)

// NewRegistry create new RegistryInterface.
func NewRegistry(clients []ClientInterface, logger logrus.FieldLogger, failFast bool) RegistryInterface {
	return &Registry{
		clients:  clients,
		logger:   logger,
		failFast: failFast,
	}
}

// Process iterate over passed plugin clients, trying to dispense and make generation request over GRPC to them.
func (r *Registry) Process(ctx context.Context, options generator.RunOptions, values generator.RunValues) (generator.RunResult, error) {
	result := make([]string, 0)
	for _, c := range r.clients {
		sub, err := c.Run(ctx, options, values)
		if err != nil {
			err = errors.Wrapf(err, `plugin "%s" error`, c.Name())
			if r.failFast {
				return nil, err
			}
			r.logger.Warn(capitaliseFirst(err.Error()))
		}
		result = append(
			result,
			funk.Map(sub, func(item string) string {
				return fmt.Sprintf(`"%s" generated: %s`, c.Name(), item)
			}).([]string)...,
		)
	}
	return result, nil
}

// HasPlugins determine is there any plugins registered.
func (r *Registry) HasPlugins() bool {
	return len(r.clients) > 0
}

func capitaliseFirst(message string) string {
	if len(message) == 0 {
		return ``
	}
	temporary := []rune(message)
	temporary[0] = unicode.ToUpper(temporary[0])
	return string(temporary)
}
