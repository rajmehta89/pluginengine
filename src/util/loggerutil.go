package util

import (
	"log"
	"sync"
	"time"
)

type Logger struct {
	instance *log.Logger
}

var (
	once sync.Once

	logInstance *Logger
)

func InitializeLogger() *Logger {

	once.Do(func() {

		logger := log.New(log.Writer(), "", log.LstdFlags|log.Lmicroseconds)

		logInstance = &Logger{instance: logger}

	})

	return logInstance

}

func (l *Logger) LogInfo(message string) {

	l.instance.Printf("[INFO] %s - %s", time.Now().Format("2006-01-02 15:04:05.000"), message)

}

func (l *Logger) LogError(err error) {

	l.instance.Printf("[ERROR] %s - %v", time.Now().Format("2006-01-02 15:04:05.000"), err)

}

func (l *Logger) LogWarning(message string) {

	l.instance.Printf("[WARNING] %s - %s", time.Now().Format("2006-01-02 15:04:05.000"), message)

}
