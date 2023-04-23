package logger

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var Trace *log.Logger
var Debug *log.Logger
var Info *log.Logger
var Warn *log.Logger
var Error *log.Logger
var Fatal *log.Logger

type LogLevel struct {
	Prefix   string
	Name     string
	AsNumber int
}

var LogLevelTrace LogLevel = LogLevel{Prefix: "TRC", Name: "trace", AsNumber: 1}
var LogLevelDebug LogLevel = LogLevel{Prefix: "DBG", Name: "debug", AsNumber: 2}
var LogLevelInfo LogLevel = LogLevel{Prefix: "INF", Name: "info", AsNumber: 3}
var LogLevelWarning LogLevel = LogLevel{Prefix: "WRN", Name: "warn", AsNumber: 4}
var LogLevelError LogLevel = LogLevel{Prefix: "ERR", Name: "error", AsNumber: 5}
var LogLevelFatal LogLevel = LogLevel{Prefix: "FTL", Name: "fatal", AsNumber: 6}

var level LogLevel = LogLevelDebug

func SetLogLevel(levelAsString string) {
	level, _ = LogLevelByString(levelAsString)
}

type logWriter struct {
	level LogLevel
}

func LogLevelByString(level string) (LogLevel, error) {
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

func LogDynDnsUpdate(service, domain, ip string, err error) {
	if err == nil {
		Info.Printf("[%v] %v -> %v (success)", service, domain, ip)
	} else {
		Error.Printf("[%v] %v -> %v (%v)", service, domain, ip, err.Error())
	}
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	if writer.level.AsNumber >= level.AsNumber {
		return fmt.Print(time.Now().UTC().Format(time.RFC3339) + " [" + writer.level.Prefix + "] " + string(bytes))
	} else {

		return fmt.Print()
	}
}

func init() {
	Trace = log.New(logWriter{level: LogLevelTrace}, "", 0)
	Debug = log.New(logWriter{level: LogLevelDebug}, "", 0)
	Info = log.New(logWriter{level: LogLevelInfo}, "", 0)
	Warn = log.New(logWriter{level: LogLevelWarning}, "", 0)
	Error = log.New(logWriter{level: LogLevelError}, "", 0)
	Fatal = log.New(logWriter{level: LogLevelInfo}, "", 0)
}
