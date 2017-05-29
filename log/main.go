// Package that provides the functionality to start und initialize the logger
package log

import (
	"log"
	"net/http"

	logger "github.com/Sirupsen/logrus"
	httpLib "github.com/genofire/golang-lib/http"
)

// Current logger with it's configuration
var Log *logger.Logger

// Function to initiate a new logger
func init() {
	Log = logger.New()
	// Enable fallback, if core logger
	log.SetOutput(Log.Writer())
}

// Function to add the information of a http request to the log
func HTTP(r *http.Request) *logger.Entry {
	return Log.WithFields(logger.Fields{
		"remote": httpLib.GetRemoteIP(r),
		"method": r.Method,
		"url":    r.URL.RequestURI(),
	})
}
