package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/ythosa/gowrapper/internal/config"
)

const componentField = "component"

func New(component string) *logrus.Logger {
	cfg := config.GetLogger()

	logger := logrus.New()
	logger.SetLevel(cfg.Level())
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
	logger.WithField(componentField, component)

	return logger
}
