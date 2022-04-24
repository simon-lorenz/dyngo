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

var LogLevel = LogLevelInfo

const (
	LogLevelTrace   = 1 // "TRC"
	LogLevelDebug   = 2 // "DBG"
	LogLevelInfo    = 3 // "INF"
	LogLevelWarning = 4 // "WRN"
	LogLevelError   = 5 // "ERR"
	LogLevelFatal   = 6 // "FTL"
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
	Trace = log.New(logWriter{level: LogLevelInfo}, "", 0)
	Debug = log.New(logWriter{level: LogLevelInfo}, "", 0)
	Info = log.New(logWriter{level: LogLevelInfo}, "", 0)
	Warn = log.New(logWriter{level: LogLevelWarning}, "", 0)
	Error = log.New(logWriter{level: LogLevelError}, "", 0)
	Fatal = log.New(logWriter{level: LogLevelInfo}, "", 0)
}
