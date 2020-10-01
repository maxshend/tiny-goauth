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

var requestDetails = Event{0, "%d %q %q %q %q"}

// RequestDetails is an HTTP request details
func (l *StandardLogger) RequestDetails(code int, method, url, agent, addr string) {
	l.Errorf(requestDetails.message, code, method, url, agent, addr)
}
