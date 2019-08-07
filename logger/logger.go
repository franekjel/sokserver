package logger

import (
	"fmt"
	"os"
	"time"
)

var exitOnError bool
var debug bool

//InitLogger set logger flags exitOnError and debug
func InitLogger(exit bool, dbg bool) {
	exitOnError = exit
	debug = dbg
}

//Debug log - in magenta, visible only with flag -d
func Debug(msg string, params ...interface{}) {
	if !debug {
		return
	}
	fmt.Print("\033[95m", "[", "DEBUG", " ", time.Now().Format("2006-01-02 15:04:05.00"), "]\033[0m ")
	fmt.Printf(msg, params...)
	fmt.Println()
}

//Info - message in green
func Info(msg string, params ...interface{}) {
	fmt.Print("\033[32m", "[", "INFO ", " ", time.Now().Format("2006-01-02 15:04:05.00"), "]\033[0m ")
	fmt.Printf(msg, params...)
	fmt.Println()
}

//Warn - message in yellow
func Warn(msg string, params ...interface{}) {
	fmt.Print("\033[93m", "[", "WARN ", " ", time.Now().Format("2006-01-02 15:04:05.00"), "]\033[0m ")
	fmt.Printf(msg, params...)
	fmt.Println()
}

//Error - message in red, exit if flag -e is set
func Error(msg string, params ...interface{}) {
	fmt.Print("\033[91m", "[", "ERROR", " ", time.Now().Format("2006-01-02 15:04:05.00"), "]\033[0m ")
	fmt.Printf(msg, params...)
	fmt.Println()
	if exitOnError {
		os.Exit(1)
	}
}

//Fatal - message in red, end program
func Fatal(msg string, params ...interface{}) {
	fmt.Print("\033[91m", "[", "FATAL", " ", time.Now().Format("2006-01-02 15:04:05.00"), "]\033[0m ")
	fmt.Printf(msg, params...)
	fmt.Println()
	os.Exit(1)
}
