package logger

import (
	"log"
	"os"
)

var Warn *log.Logger
var Info *log.Logger
var Error *log.Logger

func init() {
	Warn = log.New(os.Stderr, "[WARN] ", 0)
	Info = log.New(os.Stderr, "[INFO] ", 0)
	Error = log.New(os.Stderr, "[ERROR] ", 0)
}
