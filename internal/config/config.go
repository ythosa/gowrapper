package config

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	once sync.Once

	level logrus.Level
}

func newLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Level() logrus.Level {
	return l.level
}

func (l *Logger) init() {
	l.level = logrus.DebugLevel
}

var logger = newLogger()

func GetLogger() *Logger {
	logger.once.Do(logger.init)

	return logger
}
