package logwrapper

import "github.com/sirupsen/logrus"

// StandardLogger encapsulates logging format
type StandardLogger struct {
	*logrus.Logger
}

// Event represents a logged event
type Event struct {
	id      int
	message string
}

// New initialize standard logger
func New() *StandardLogger {
	baseLogger := logrus.New()
	standardLogger := &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}

	return standardLogger
}
