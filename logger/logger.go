package logger

import (
	"fmt"
	"os"
	"time"
)

var exitOnError bool
var debug bool

const (
	//DEBUG - additional debugging information. Magenta
	DEBUG = "DEBUG"
	//INFO - standard information, Green
	INFO = "INFO "
	//WARN - warning, there is problem, but it shouldn't affect stability. Yellow
	WARN = "WARN "
	//ERR - there is error and SOK may be unstable. Red
	ERR = "ERROR"
	//FATAL - there is error and SOK going to exit. Red
	FATAL = "FATAL"
)

//InitLogger set logger flags exitOnError and debug
func InitLogger(exit bool, dbg bool) {
	exitOnError = exit
	debug = dbg
}

func logLevelToColor(level string) string {

	switch level {
	case DEBUG:
		return "\033[95m"
	case INFO:
		return "\033[32m"
	case WARN:
		return "\033[93m"
	case ERR:
		return "\033[91m"
	case FATAL:
		return "\033[31m"
	}

	return "\033[96m"
}

//Log printed to stdout. logLevel specify color and behaviur: DEBUG - magenta (only if debug is true), INFO - green, WARN - yellow, ERR and FATAL -red. FATAL also end program
func Log(logLevel, msg string, params ...interface{}) {
	if logLevel == DEBUG && !debug {
		return
	}

	fmt.Print(logLevelToColor(logLevel), "[", logLevel, " ", time.Now().Format("2006-01-02 15:04:05.00"), "]\033[0m ")
	fmt.Printf(msg, params...)
	fmt.Println()

	if logLevel == FATAL || (logLevel == ERR && exitOnError) {
		os.Exit(1)
	}
}
