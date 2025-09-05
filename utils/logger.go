package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
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

// Logger interface defines the logging methods
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	SetLevel(level LogLevel)
	SetOutput(writer io.Writer)
}

// ModuleLogger wraps a logger with module name
type ModuleLogger struct {
	logger Logger
	module string
}

// NewModuleLogger creates a new module logger
func NewModuleLogger(logger Logger, module string) *ModuleLogger {
	return &ModuleLogger{
		logger: logger,
		module: module,
	}
}

// Debug logs a debug message with module name
func (ml *ModuleLogger) Debug(format string, args ...interface{}) {
	if sl, ok := ml.logger.(*SimpleLogger); ok {
		sl.logWithModule(DEBUG, ml.module, format, args...)
	} else {
		ml.logger.Debug(format, args...)
	}
}

// Info logs an info message with module name
func (ml *ModuleLogger) Info(format string, args ...interface{}) {
	if sl, ok := ml.logger.(*SimpleLogger); ok {
		sl.logWithModule(INFO, ml.module, format, args...)
	} else {
		ml.logger.Info(format, args...)
	}
}

// Warn logs a warning message with module name
func (ml *ModuleLogger) Warn(format string, args ...interface{}) {
	if sl, ok := ml.logger.(*SimpleLogger); ok {
		sl.logWithModule(WARN, ml.module, format, args...)
	} else {
		ml.logger.Warn(format, args...)
	}
}

// Error logs an error message with module name
func (ml *ModuleLogger) Error(format string, args ...interface{}) {
	if sl, ok := ml.logger.(*SimpleLogger); ok {
		sl.logWithModule(ERROR, ml.module, format, args...)
	} else {
		ml.logger.Error(format, args...)
	}
}

// Fatal logs a fatal message with module name
func (ml *ModuleLogger) Fatal(format string, args ...interface{}) {
	if sl, ok := ml.logger.(*SimpleLogger); ok {
		sl.logWithModule(FATAL, ml.module, format, args...)
	} else {
		ml.logger.Fatal(format, args...)
	}
}

// SetLevel sets the logging level
func (ml *ModuleLogger) SetLevel(level LogLevel) {
	ml.logger.SetLevel(level)
}

// SetOutput sets the output writer
func (ml *ModuleLogger) SetOutput(writer io.Writer) {
	ml.logger.SetOutput(writer)
}

// SimpleLogger implements the Logger interface
type SimpleLogger struct {
	level  LogLevel
	output io.Writer
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel) Logger {
	return &SimpleLogger{
		level:  level,
		output: os.Stdout,
	}
}

// NewLoggerWithOutput creates a new logger with custom output
func NewLoggerWithOutput(level LogLevel, output io.Writer) Logger {
	return &SimpleLogger{
		level:  level,
		output: output,
	}
}

// SetLevel sets the logging level
func (l *SimpleLogger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput sets the output writer
func (l *SimpleLogger) SetOutput(writer io.Writer) {
	l.output = writer
}

// log writes a log message with the specified level
func (l *SimpleLogger) log(level LogLevel, format string, args ...interface{}) {
	l.logWithModule(level, "", format, args...)
}

// logWithModule writes a log message with the specified level and module
func (l *SimpleLogger) logWithModule(level LogLevel, module, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)

	// Color codes for different log levels
	var colorCode string
	switch level {
	case DEBUG:
		colorCode = "\033[36m" // Cyan
	case INFO:
		colorCode = "\033[32m" // Green
	case WARN:
		colorCode = "\033[33m" // Yellow
	case ERROR:
		colorCode = "\033[31m" // Red
	case FATAL:
		colorCode = "\033[35m" // Magenta
	}

	resetCode := "\033[0m"

	// Format log line with optional module name
	var logLine string
	if module != "" {
		logLine = fmt.Sprintf("[%s] %s%s%s [%s] %s\n",
			timestamp,
			colorCode,
			level.String(),
			resetCode,
			module,
			message)
	} else {
		logLine = fmt.Sprintf("[%s] %s%s%s %s\n",
			timestamp,
			colorCode,
			level.String(),
			resetCode,
			message)
	}

	l.output.Write([]byte(logLine))
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *SimpleLogger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *SimpleLogger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *SimpleLogger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *SimpleLogger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// ParseLogLevel parses a string to LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO // Default to INFO level
	}
}

// Global logger instance
var globalLogger Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(level LogLevel) {
	globalLogger = NewLogger(level)
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() Logger {
	if globalLogger == nil {
		globalLogger = NewLogger(INFO)
	}
	return globalLogger
}

// Convenience functions for global logger
func Debug(format string, args ...interface{}) {
	GetGlobalLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetGlobalLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetGlobalLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetGlobalLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetGlobalLogger().Fatal(format, args...)
}

// GetModuleLogger returns a module-specific logger
func GetModuleLogger(module string) *ModuleLogger {
	return NewModuleLogger(GetGlobalLogger(), module)
}
