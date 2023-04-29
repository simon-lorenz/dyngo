package logger

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// Global application logger
var Trace *log.Logger
var Debug *log.Logger
var Info *log.Logger
var Warn *log.Logger
var Error *log.Logger
var Fatal *log.Logger

// Service logger
type LoggerCollection struct {
	Trace *log.Logger
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Fatal *log.Logger
}

type logLevel struct {
	Prefix   string
	Name     string
	Priority int
}

type logWriter struct {
	level logLevel
}

var LogLevelTrace logLevel = logLevel{Prefix: "TRC", Name: "trace", Priority: 1}
var LogLevelDebug logLevel = logLevel{Prefix: "DBG", Name: "debug", Priority: 2}
var LogLevelInfo logLevel = logLevel{Prefix: "INF", Name: "info", Priority: 3}
var LogLevelWarning logLevel = logLevel{Prefix: "WRN", Name: "warn", Priority: 4}
var LogLevelError logLevel = logLevel{Prefix: "ERR", Name: "error", Priority: 5}
var LogLevelFatal logLevel = logLevel{Prefix: "FTL", Name: "fatal", Priority: 6}

var level logLevel = LogLevelDebug

func SetLogLevel(levelAsString string) {
	level, _ = getLogLevelByString(levelAsString)
}

func getLogLevelByString(level string) (logLevel, error) {
	if level == LogLevelTrace.Name {
		return LogLevelTrace, nil
	} else if level == LogLevelDebug.Name {
		return LogLevelDebug, nil
	} else if level == LogLevelInfo.Name {
		return LogLevelInfo, nil
	} else if level == LogLevelWarning.Name {
		return LogLevelWarning, nil
	} else if level == LogLevelError.Name {
		return LogLevelError, nil
	} else if level == LogLevelFatal.Name {
		return LogLevelFatal, nil
	} else {
		return LogLevelInfo, errors.New("Cannot determine log level for string \"" + level + "\"")
	}
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	if writer.level.Priority >= level.Priority {
		return fmt.Print(time.Now().UTC().Format(time.RFC3339) + " [" + writer.level.Prefix + "] " + string(bytes))
	} else {

		return fmt.Print()
	}
}

func NewServiceLoggerCollection(serviceName string) *LoggerCollection {
	return &LoggerCollection{
		Trace: log.New(logWriter{level: LogLevelTrace}, "["+serviceName+"] ", 0),
		Debug: log.New(logWriter{level: LogLevelDebug}, "["+serviceName+"] ", 0),
		Info:  log.New(logWriter{level: LogLevelInfo}, "["+serviceName+"] ", 0),
		Warn:  log.New(logWriter{level: LogLevelWarning}, "["+serviceName+"] ", 0),
		Error: log.New(logWriter{level: LogLevelError}, "["+serviceName+"] ", 0),
		Fatal: log.New(logWriter{level: LogLevelFatal}, "["+serviceName+"] ", 0),
	}

}

func init() {
	Trace = log.New(logWriter{level: LogLevelTrace}, "[app] ", 0)
	Debug = log.New(logWriter{level: LogLevelDebug}, "[app] ", 0)
	Info = log.New(logWriter{level: LogLevelInfo}, "[app] ", 0)
	Warn = log.New(logWriter{level: LogLevelWarning}, "[app] ", 0)
	Error = log.New(logWriter{level: LogLevelError}, "[app] ", 0)
	Fatal = log.New(logWriter{level: LogLevelFatal}, "[app] ", 0)
}
