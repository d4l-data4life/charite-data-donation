package logging

import (
	"fmt"
	"log"
	"sync"
)

var once sync.Once
var instance *log.Logger

func getLogger() *log.Logger {
	once.Do(func() {
		instance = log.Default()
	})
	return instance
}

// LogError logs an error with the singleton logger with message and error
func LogError(message string, err error) {
	getLogger().Printf("[Error]: %s - %s", message, err.Error())
}

// LogWarning logs a warning with the singleton logger with message and error
func LogWarning(message string, err error) {
	getLogger().Printf("[Warn]: %s - %s", message, err.Error())
}

// LogInfo logs an info with the singleton logger with message and error
func LogInfo(message string) {
	getLogger().Printf("[Info]: %s", message)
}

// LogInfof info-level log with formatting
func LogInfof(format string, fields ...interface{}) {
	msg := fmt.Sprintf(format, fields...)
	getLogger().Printf("[Info]: %s", msg)
}
