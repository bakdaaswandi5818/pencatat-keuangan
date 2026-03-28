package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New returns a configured logrus logger that writes JSON to stdout.
func New() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	return log
}
