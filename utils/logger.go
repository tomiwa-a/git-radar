package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logger *log.Logger
var logFile *os.File

func InitLogger() error {
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(logsDir, "performance.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	logFile = f
	logger = log.New(f, "", log.LstdFlags|log.Lmicroseconds)
	logger.Println("=== Git-Radar started ===")
	return nil
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func LogTiming(operation string, duration time.Duration) {
	if logger != nil {
		logger.Printf("[PERF] %s: %v", operation, duration)
	}
}

func LogEvent(event string) {
	if logger != nil {
		logger.Printf("[EVENT] %s", event)
	}
}

func TimeOperation(name string, fn func()) {
	start := time.Now()
	fn()
	LogTiming(name, time.Since(start))
}

func TimedFunc[T any](name string, fn func() T) T {
	start := time.Now()
	result := fn()
	LogTiming(name, time.Since(start))
	return result
}

func TimedFuncErr[T any](name string, fn func() (T, error)) (T, error) {
	start := time.Now()
	result, err := fn()
	if err != nil {
		LogTiming(fmt.Sprintf("%s (error: %v)", name, err), time.Since(start))
	} else {
		LogTiming(name, time.Since(start))
	}
	return result, err
}
