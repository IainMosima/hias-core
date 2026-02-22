package utils

import (
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	WarnLog  *log.Logger
	ErrorLog *log.Logger
	DebugLog *log.Logger
)

func init() {
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLog = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(format string, v ...interface{}) {
	InfoLog.Printf(format, v...)
}

func LogWarn(format string, v ...interface{}) {
	WarnLog.Printf(format, v...)
}

func LogError(format string, v ...interface{}) {
	ErrorLog.Printf(format, v...)
}

func LogDebug(format string, v ...interface{}) {
	DebugLog.Printf(format, v...)
}
