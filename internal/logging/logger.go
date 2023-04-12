package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/ythosa/gowrapper/internal/config"
)

const componentField = "component"

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

func New(component string) Logger {
	cfg := config.GetLogger()

	logger := logrus.New()
	logger.SetLevel(cfg.Level())
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})

	return logger.WithFields(logrus.Fields{componentField: component})
}
