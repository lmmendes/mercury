package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	*log.Logger
	level Level
}

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func New(out io.Writer, level Level) *Logger {
	return &Logger{
		Logger: log.New(out, "", 0),
		level:  level,
	}
}

func (l *Logger) log(level Level, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	// Extract just the filename from the full path
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		file = file[idx+1:]
	}

	// Format the message
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logLine := fmt.Sprintf("[%s] %-5s %s:%d: %s", timestamp, level, file, line, msg)

	l.Logger.Output(2, logLine)

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FATAL, format, v...)
}

// ErrorWithStack logs an error with its stack trace
func (l *Logger) ErrorWithStack(err error) {
	if err == nil {
		return
	}

	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])

	l.Error("Error: %v\nStack Trace:\n%s", err, stackTrace)
}

// Add these methods to the Level type
func (l *Level) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "debug":
		*l = DEBUG
	case "info":
		*l = INFO
	case "warn":
		*l = WARN
	case "error":
		*l = ERROR
	case "fatal":
		*l = FATAL
	default:
		return fmt.Errorf("unknown log level: %s", text)
	}
	return nil
}

func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
