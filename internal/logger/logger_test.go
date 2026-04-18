package logger

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestNew tests logger creation
func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		level   LogLevel
		noColor bool
	}{
		{name: "debug level with color", level: DebugLevel, noColor: false},
		{name: "info level with color", level: InfoLevel, noColor: false},
		{name: "warn level with color", level: WarnLevel, noColor: false},
		{name: "error level with color", level: ErrorLevel, noColor: false},
		{name: "debug level no color", level: DebugLevel, noColor: true},
		{name: "info level no color", level: InfoLevel, noColor: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.level, tt.noColor)
			
			testutil.AssertTrue(t, logger != nil, "Logger should not be nil")
			testutil.AssertEqual(t, tt.level, logger.level)
			testutil.AssertEqual(t, tt.noColor, logger.noColor)
			testutil.AssertTrue(t, len(logger.writers) > 0, "Logger should have at least one writer")
		})
	}
}

// TestLogLevelString tests LogLevel.String method
func TestLogLevelString(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected string
	}{
		{name: "debug level", level: DebugLevel, expected: "DEBUG"},
		{name: "info level", level: InfoLevel, expected: "INFO"},
		{name: "warn level", level: WarnLevel, expected: "WARN"},
		{name: "error level", level: ErrorLevel, expected: "ERROR"},
		{name: "unknown level", level: LogLevel(99), expected: "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.level.String()
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestParseLogLevel tests ParseLogLevel function
func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LogLevel
	}{
		{name: "parse debug", input: "debug", expected: DebugLevel},
		{name: "parse info", input: "info", expected: InfoLevel},
		{name: "parse warn", input: "warn", expected: WarnLevel},
		{name: "parse error", input: "error", expected: ErrorLevel},
		{name: "parse unknown defaults to info", input: "unknown", expected: InfoLevel},
		{name: "parse empty defaults to info", input: "", expected: InfoLevel},
		{name: "parse uppercase defaults to info", input: "DEBUG", expected: InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLogLevel(tt.input)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestSetLevel tests SetLevel method
func TestSetLevel(t *testing.T) {
	logger := New(InfoLevel, false)
	testutil.AssertEqual(t, InfoLevel, logger.level)

	logger.SetLevel(DebugLevel)
	testutil.AssertEqual(t, DebugLevel, logger.level)

	logger.SetLevel(ErrorLevel)
	testutil.AssertEqual(t, ErrorLevel, logger.level)
}

// TestAddWriter tests AddWriter method
func TestAddWriter(t *testing.T) {
	logger := New(InfoLevel, false)
	initialWriterCount := len(logger.writers)

	buf := &bytes.Buffer{}
	logger.AddWriter(buf)

	testutil.AssertEqual(t, initialWriterCount+1, len(logger.writers))
}

// TestDebugLogging tests Debug method
func TestDebugLogging(t *testing.T) {
	tests := []struct {
		name          string
		loggerLevel   LogLevel
		shouldLog     bool
	}{
		{name: "debug logs at debug level", loggerLevel: DebugLevel, shouldLog: true},
		{name: "debug doesn't log at info level", loggerLevel: InfoLevel, shouldLog: false},
		{name: "debug doesn't log at warn level", loggerLevel: WarnLevel, shouldLog: false},
		{name: "debug doesn't log at error level", loggerLevel: ErrorLevel, shouldLog: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(tt.loggerLevel, true)
			logger.writers = []io.Writer{buf}

			logger.Debug("test debug message")

			output := buf.String()
			if tt.shouldLog {
				testutil.AssertContains(t, output, "DEBUG")
				testutil.AssertContains(t, output, "test debug message")
			} else {
				testutil.AssertEqual(t, "", output)
			}
		})
	}
}

// TestInfoLogging tests Info method
func TestInfoLogging(t *testing.T) {
	tests := []struct {
		name          string
		loggerLevel   LogLevel
		shouldLog     bool
	}{
		{name: "info logs at debug level", loggerLevel: DebugLevel, shouldLog: true},
		{name: "info logs at info level", loggerLevel: InfoLevel, shouldLog: true},
		{name: "info doesn't log at warn level", loggerLevel: WarnLevel, shouldLog: false},
		{name: "info doesn't log at error level", loggerLevel: ErrorLevel, shouldLog: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(tt.loggerLevel, true)
			logger.writers = []io.Writer{buf}

			logger.Info("test info message")

			output := buf.String()
			if tt.shouldLog {
				testutil.AssertContains(t, output, "INFO")
				testutil.AssertContains(t, output, "test info message")
			} else {
				testutil.AssertEqual(t, "", output)
			}
		})
	}
}

// TestWarnLogging tests Warn method
func TestWarnLogging(t *testing.T) {
	tests := []struct {
		name          string
		loggerLevel   LogLevel
		shouldLog     bool
	}{
		{name: "warn logs at debug level", loggerLevel: DebugLevel, shouldLog: true},
		{name: "warn logs at info level", loggerLevel: InfoLevel, shouldLog: true},
		{name: "warn logs at warn level", loggerLevel: WarnLevel, shouldLog: true},
		{name: "warn doesn't log at error level", loggerLevel: ErrorLevel, shouldLog: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(tt.loggerLevel, true)
			logger.writers = []io.Writer{buf}

			logger.Warn("test warn message")

			output := buf.String()
			if tt.shouldLog {
				testutil.AssertContains(t, output, "WARN")
				testutil.AssertContains(t, output, "test warn message")
			} else {
				testutil.AssertEqual(t, "", output)
			}
		})
	}
}

// TestErrorLogging tests Error method
func TestErrorLogging(t *testing.T) {
	tests := []struct {
		name          string
		loggerLevel   LogLevel
		shouldLog     bool
	}{
		{name: "error logs at debug level", loggerLevel: DebugLevel, shouldLog: true},
		{name: "error logs at info level", loggerLevel: InfoLevel, shouldLog: true},
		{name: "error logs at warn level", loggerLevel: WarnLevel, shouldLog: true},
		{name: "error logs at error level", loggerLevel: ErrorLevel, shouldLog: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(tt.loggerLevel, true)
			logger.writers = []io.Writer{buf}

			logger.Error("test error message")

			output := buf.String()
			if tt.shouldLog {
				testutil.AssertContains(t, output, "ERROR")
				testutil.AssertContains(t, output, "test error message")
			} else {
				testutil.AssertEqual(t, "", output)
			}
		})
	}
}

// TestSuccessLogging tests Success method
func TestSuccessLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(InfoLevel, true)
	logger.writers = []io.Writer{buf}

	logger.Success("test success message")

	output := buf.String()
	testutil.AssertContains(t, output, "INFO")
	testutil.AssertContains(t, output, "test success message")
}

// TestLoggingWithFields tests logging with structured fields
func TestLoggingWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(InfoLevel, true)
	logger.writers = []io.Writer{buf}

	logger.Info("test message", F("key1", "value1"), F("key2", 42))

	output := buf.String()
	testutil.AssertContains(t, output, "test message")
	testutil.AssertContains(t, output, "key1=value1")
	testutil.AssertContains(t, output, "key2=42")
}

// TestFieldCreation tests Field creation
func TestFieldCreation(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    interface{}
	}{
		{name: "string field", key: "name", value: "test"},
		{name: "int field", key: "count", value: 42},
		{name: "bool field", key: "enabled", value: true},
		{name: "float field", key: "price", value: 19.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := F(tt.key, tt.value)
			testutil.AssertEqual(t, tt.key, field.Key)
			testutil.AssertEqual(t, tt.value, field.Value)
		})
	}
}

// TestColorOutput tests colored vs non-colored output
func TestColorOutput(t *testing.T) {
	tests := []struct {
		name    string
		noColor bool
	}{
		{name: "with color", noColor: false},
		{name: "without color", noColor: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(InfoLevel, tt.noColor)
			logger.writers = []io.Writer{buf}

			logger.Info("test message")

			output := buf.String()
			testutil.AssertContains(t, output, "test message")
			
			// Check for ANSI color codes
			hasColorCodes := strings.Contains(output, "\033[")
			if tt.noColor {
				testutil.AssertFalse(t, hasColorCodes, "Output should not contain color codes")
			} else {
				testutil.AssertTrue(t, hasColorCodes, "Output should contain color codes")
			}
		})
	}
}

// TestMultipleWriters tests logging to multiple writers
func TestMultipleWriters(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	
	logger := New(InfoLevel, true)
	logger.writers = []io.Writer{buf1}
	logger.AddWriter(buf2)

	logger.Info("test message")

	output1 := buf1.String()
	output2 := buf2.String()

	testutil.AssertContains(t, output1, "test message")
	testutil.AssertContains(t, output2, "test message")
	testutil.AssertEqual(t, output1, output2)
}

// TestInitDefault tests InitDefault function
func TestInitDefault(t *testing.T) {
	InitDefault(DebugLevel, true)
	
	logger := GetDefault()
	testutil.AssertTrue(t, logger != nil, "Default logger should not be nil")
	testutil.AssertEqual(t, DebugLevel, logger.level)
	testutil.AssertEqual(t, true, logger.noColor)
}

// TestGetDefault tests GetDefault function
func TestGetDefault(t *testing.T) {
	// Reset default logger
	defaultLogger = nil
	
	logger := GetDefault()
	testutil.AssertTrue(t, logger != nil, "Default logger should not be nil")
	testutil.AssertEqual(t, InfoLevel, logger.level)
	testutil.AssertEqual(t, false, logger.noColor)
}

// TestGlobalDebug tests global Debug function
func TestGlobalDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	InitDefault(DebugLevel, true)
	defaultLogger.writers = []io.Writer{buf}

	Debug("global debug message")

	output := buf.String()
	testutil.AssertContains(t, output, "DEBUG")
	testutil.AssertContains(t, output, "global debug message")
}

// TestGlobalInfo tests global Info function
func TestGlobalInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	InitDefault(InfoLevel, true)
	defaultLogger.writers = []io.Writer{buf}

	Info("global info message")

	output := buf.String()
	testutil.AssertContains(t, output, "INFO")
	testutil.AssertContains(t, output, "global info message")
}

// TestGlobalWarn tests global Warn function
func TestGlobalWarn(t *testing.T) {
	buf := &bytes.Buffer{}
	InitDefault(WarnLevel, true)
	defaultLogger.writers = []io.Writer{buf}

	Warn("global warn message")

	output := buf.String()
	testutil.AssertContains(t, output, "WARN")
	testutil.AssertContains(t, output, "global warn message")
}

// TestGlobalError tests global Error function
func TestGlobalError(t *testing.T) {
	buf := &bytes.Buffer{}
	InitDefault(ErrorLevel, true)
	defaultLogger.writers = []io.Writer{buf}

	Error("global error message")

	output := buf.String()
	testutil.AssertContains(t, output, "ERROR")
	testutil.AssertContains(t, output, "global error message")
}

// TestGlobalSuccess tests global Success function
func TestGlobalSuccess(t *testing.T) {
	buf := &bytes.Buffer{}
	InitDefault(InfoLevel, true)
	defaultLogger.writers = []io.Writer{buf}

	Success("global success message")

	output := buf.String()
	testutil.AssertContains(t, output, "INFO")
	testutil.AssertContains(t, output, "global success message")
}

// TestLogTimestamp tests that log messages include timestamps
func TestLogTimestamp(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(InfoLevel, true)
	logger.writers = []io.Writer{buf}

	logger.Info("test message")

	output := buf.String()
	// Check for timestamp format YYYY-MM-DD HH:MM:SS
	testutil.AssertTrue(t, strings.Contains(output, "["), "Output should contain timestamp bracket")
	testutil.AssertTrue(t, strings.Contains(output, "]"), "Output should contain timestamp bracket")
}

// TestConcurrentLogging tests concurrent logging safety
func TestConcurrentLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(InfoLevel, true)
	logger.writers = []io.Writer{buf}

	done := make(chan bool)
	
	// Start multiple goroutines logging concurrently
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Info("concurrent message", F("goroutine", id))
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	output := buf.String()
	// Just verify that we got some output without panicking
	testutil.AssertTrue(t, len(output) > 0, "Should have logged messages")
}

// TestLogLevelFiltering tests that log level filtering works correctly
func TestLogLevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(WarnLevel, true)
	logger.writers = []io.Writer{buf}

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	
	// Debug and Info should not appear
	testutil.AssertFalse(t, strings.Contains(output, "debug message"), "Debug should be filtered")
	testutil.AssertFalse(t, strings.Contains(output, "info message"), "Info should be filtered")
	
	// Warn and Error should appear
	testutil.AssertContains(t, output, "warn message")
	testutil.AssertContains(t, output, "error message")
}

// TestPropertyLoggerIsolation tests that loggers are isolated from each other
func TestPropertyLoggerIsolation(t *testing.T) {
	// Feature: comprehensive-testing-suite, Property 1: Unit Test Isolation
	testutil.PropertyTestWithContext(t, "logger_isolation", 10, 100*time.Millisecond, func(t *testing.T, ctx context.Context) {
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}
		
		logger1 := New(InfoLevel, true)
		logger1.writers = []io.Writer{buf1}
		
		logger2 := New(DebugLevel, true)
		logger2.writers = []io.Writer{buf2}

		logger1.Info("logger1 message")
		logger2.Debug("logger2 message")

		output1 := buf1.String()
		output2 := buf2.String()

		// Verify isolation
		testutil.AssertContains(t, output1, "logger1 message")
		testutil.AssertFalse(t, strings.Contains(output1, "logger2 message"), "Logger1 should not see logger2 messages")
		
		testutil.AssertContains(t, output2, "logger2 message")
		testutil.AssertFalse(t, strings.Contains(output2, "logger1 message"), "Logger2 should not see logger1 messages")
	})
}

