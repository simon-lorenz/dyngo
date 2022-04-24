package logger

import (
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

var LogLevel = logLevelInfo

const (
	logLevelTrace   = 1 // "TRC"
	logLevelDebug   = 2 // "DBG"
	logLevelInfo    = 3 // "INF"
	logLevelWarning = 4 // "WRN"
	logLevelError   = 5 // "ERR"
	logLevelFatal   = 6 // "FTL"
)

type logWriter struct {
	level int
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	var level string

	switch writer.level {
	case 1:
		level = "TRC"
		break
	case 2:
		level = "DBG"
		break
	case 3:
		level = "INF"
		break
	case 4:
		level = "WRN"
		break
	case 5:
		level = "ERR"
		break
	case 6:
		level = "FTL"
		break
	}

	if writer.level >= LogLevel {
		return fmt.Print(time.Now().UTC().Format(time.RFC3339) + " [" + level + "] " + string(bytes))
	} else {

		return fmt.Print()
	}
}

func init() {
	Trace = log.New(logWriter{level: logLevelInfo}, "", 0)
	Debug = log.New(logWriter{level: logLevelInfo}, "", 0)
	Info = log.New(logWriter{level: logLevelInfo}, "", 0)
	Warn = log.New(logWriter{level: logLevelWarning}, "", 0)
	Error = log.New(logWriter{level: logLevelError}, "", 0)
	Fatal = log.New(logWriter{level: logLevelInfo}, "", 0)
}
