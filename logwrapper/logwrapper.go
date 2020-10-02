package logwrapper

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

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

var requestDetails = Event{0, "%d %s %s %s %s"}
var requestError = Event{1, "%s %s %s %s caused %q"}
var fatalError = Event{2, "Application stopped: %s"}

// RequestDetails logs an HTTP request details
func (l *StandardLogger) RequestDetails(r *http.Request, code int) {
	l.Errorf(requestDetails.message, code, r.Method, r.RequestURI, r.UserAgent(), r.RemoteAddr)
}

// RequestError logs errors that come up while request handling
func (l *StandardLogger) RequestError(r *http.Request, err error) {
	l.Errorf(requestError.message, r.Method, r.RequestURI, r.UserAgent(), r.RemoteAddr, err)
}

// FatalError logs about fatal errors
func (l *StandardLogger) FatalError(err error) {
	l.Fatalf(fatalError.message, err)
}
