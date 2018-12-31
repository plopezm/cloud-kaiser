package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger - Contains the standard golang log and the log level set.
type Logger struct {
	log   *log.Logger
	level uint8
}

// Info - Prints using INFO logger level
func (logger *Logger) Info(v ...interface{}) {
	if logger.level >= 3 {
		logger.log.Output(2, fmt.Sprintf("[INFO] %s\n", v...))
	}
}

// Warn - Prints using WARN logger level
func (logger *Logger) Warn(v ...interface{}) {
	if logger.level >= 2 {
		logger.log.Output(2, fmt.Sprintf("[WARN] %s\n", v...))
	}
}

// Error - Prints using ERROR logger level
func (logger *Logger) Error(v ...interface{}) {
	logger.log.Output(2, fmt.Sprintf("[ERROR] %s\n", v...))
}

// Fatal - Prints and then finishes the execution
func (logger *Logger) Fatal(v ...interface{}) {
	logger.log.Output(2, fmt.Sprintf("[FATAL] %s\n", v...))
	os.Exit(1)
}

// Debug - Prints using DEBUG logger level
func (logger *Logger) Debug(v ...interface{}) {
	if logger.level >= 4 {
		logger.log.Output(2, fmt.Sprintf("[DEBUG] %s\n", v...))
	}
}

var logger *Logger

func InitializeLogger(defaultLogLevel string) {
	logger = new(Logger)
	logger.log = log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Ltime|log.LUTC)
	switch defaultLogLevel {
	case "DEBUG":
		logger.level = 4
	case "INFO":
		logger.level = 3
	case "WARN":
		logger.level = 2
	default:
		logger.level = 1

	}
}

// GetLogger - Returns the singleton logger with the level set from config file
func GetLogger() *Logger {
	if logger == nil {
		panic("Logger not initalized, please use InitializeLogger first")
	}
	return logger
}
