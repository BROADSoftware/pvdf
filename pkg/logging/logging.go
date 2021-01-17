package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var logLevelByString = map[string]logrus.Level{
	"PANIC": logrus.PanicLevel,
	"FATAL": logrus.FatalLevel,
	"ERROR": logrus.ErrorLevel,
	"WARN":  logrus.WarnLevel,
	"INFO":  logrus.InfoLevel,
	"DEBUG": logrus.DebugLevel,
	"TRACE": logrus.TraceLevel,
}

var Log *logrus.Logger

func ConfigLogger(level string, json bool) {
	llevel := logLevelByString[strings.ToUpper(level)]
	if llevel == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "\nInvalid -logLevel value '%s'. Must be one of PANIC, FATAL, WARNING, WARN, INFO, DEBUG or TRACE\n", level)
		os.Exit(2)
	}
	Log.SetLevel(llevel)
	if json {
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{})
	}
}

func init() {
	Log = logrus.New()
	Log.Out = os.Stdout
	ConfigLogger("INFO", true)
}
