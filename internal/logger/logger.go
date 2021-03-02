package logger

import (
	"github.com/sirupsen/logrus"
)

var (
	fieldLogger *logrus.Logger
)

func init() {
	fieldLogger = logrus.New()
}

// GetLogger return internal predefined logrus.FieldLogger interface implementation.
func GetLogger() logrus.FieldLogger {
	return fieldLogger
}
