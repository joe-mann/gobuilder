package view

import "github.com/sirupsen/logrus"

func Always(logger *logrus.Logger, format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Verbose(logger *logrus.Logger, format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
