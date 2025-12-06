package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging capabilities
type Logger struct {
	level    LogLevel
	noColor  bool
	writers  []io.Writer
}

// New creates a new logger instance
func New(level LogLevel, noColor bool) *Logger {
	return &Logger{
		level:   level,
		noColor: noColor,
		writers: []io.Writer{os.Stdout},
	}
}

// AddWriter adds an additional output writer
func (l *Logger) AddWriter(w io.Writer) {
	l.writers = append(l.writers, w)
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...Field) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, msg, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...Field) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, msg, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...Field) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, msg, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...Field) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, msg, fields...)
	}
}

// Success logs a success message (info level with green color)
func (l *Logger) Success(msg string, fields ...Field) {
	if l.level <= InfoLevel {
		l.logColored(InfoLevel, msg, ColorGreen, fields...)
	}
}

// log writes a log message to all writers
func (l *Logger) log(level LogLevel, msg string, fields ...Field) {
	color := l.getColor(level)
	l.logColored(level, msg, color, fields...)
}

// logColored writes a colored log message
func (l *Logger) logColored(level LogLevel, msg string, color Color, fields ...Field) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := level.String()
	
	// Format the message
	var output string
	if l.noColor {
		output = fmt.Sprintf("[%s] %s: %s", timestamp, levelStr, msg)
	} else {
		coloredLevel := Colorize(levelStr, color)
		output = fmt.Sprintf("[%s] %s: %s", timestamp, coloredLevel, msg)
	}
	
	// Add fields if any
	if len(fields) > 0 {
		output += " |"
		for _, field := range fields {
			output += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}
	
	output += "\n"
	
	// Write to all writers
	for _, w := range l.writers {
		fmt.Fprint(w, output)
	}
}

// getColor returns the color for a log level
func (l *Logger) getColor(level LogLevel) Color {
	switch level {
	case DebugLevel:
		return ColorCyan
	case InfoLevel:
		return ColorBlue
	case WarnLevel:
		return ColorYellow
	case ErrorLevel:
		return ColorRed
	default:
		return ColorReset
	}
}

// Field represents a structured log field
type Field struct {
	Key   string
	Value interface{}
}

// F creates a new field
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// Global logger instance
var defaultLogger *Logger

// InitDefault initializes the default logger
func InitDefault(level LogLevel, noColor bool) {
	defaultLogger = New(level, noColor)
}

// GetDefault returns the default logger
func GetDefault() *Logger {
	if defaultLogger == nil {
		defaultLogger = New(InfoLevel, false)
	}
	return defaultLogger
}

// Convenience functions using the default logger

// Debug logs a debug message using the default logger
func Debug(msg string, fields ...Field) {
	GetDefault().Debug(msg, fields...)
}

// Info logs an info message using the default logger
func Info(msg string, fields ...Field) {
	GetDefault().Info(msg, fields...)
}

// Warn logs a warning message using the default logger
func Warn(msg string, fields ...Field) {
	GetDefault().Warn(msg, fields...)
}

// Error logs an error message using the default logger
func Error(msg string, fields ...Field) {
	GetDefault().Error(msg, fields...)
}

// Success logs a success message using the default logger
func Success(msg string, fields ...Field) {
	GetDefault().Success(msg, fields...)
}
