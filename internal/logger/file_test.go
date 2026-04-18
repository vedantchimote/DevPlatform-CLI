package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestNewFileLogger tests creating a new file logger
func TestNewFileLogger(t *testing.T) {
	fl, err := NewFileLogger()
	testutil.AssertNoError(t, err)
	testutil.AssertTrue(t, fl != nil, "FileLogger should not be nil")
	
	// Clean up
	defer fl.Close()
	
	// Verify log directory was created
	testutil.AssertTrue(t, fl.logDir != "", "Log directory should be set")
	
	// Verify log file was created
	logFile := fl.GetCurrentLogFile()
	testutil.AssertTrue(t, logFile != "", "Log file should be set")
	testutil.AssertTrue(t, strings.Contains(logFile, "devplatform_"), "Log file should have correct prefix")
	testutil.AssertTrue(t, strings.HasSuffix(logFile, ".log"), "Log file should have .log extension")
}

// TestFileLoggerWrite tests writing to the file logger
func TestFileLoggerWrite(t *testing.T) {
	fl, err := NewFileLogger()
	testutil.AssertNoError(t, err)
	defer fl.Close()
	
	// Write some data
	testData := "test log message\n"
	n, err := fl.Write([]byte(testData))
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, len(testData), n)
	
	// Close to flush
	fl.Close()
	
	// Read back and verify
	content, err := os.ReadFile(fl.GetCurrentLogFile())
	testutil.AssertNoError(t, err)
	testutil.AssertContains(t, string(content), "test log message")
}

// TestFileLoggerWriteJSON tests writing JSON log entries
func TestFileLoggerWriteJSON(t *testing.T) {
	fl, err := NewFileLogger()
	testutil.AssertNoError(t, err)
	defer fl.Close()
	
	// Write JSON entry
	fields := map[string]interface{}{
		"user":   "testuser",
		"action": "login",
		"status": "success",
	}
	
	err = fl.WriteJSON("INFO", "User logged in", fields)
	testutil.AssertNoError(t, err)
	
	// Close to flush
	fl.Close()
	
	// Read back and verify
	content, err := os.ReadFile(fl.GetCurrentLogFile())
	testutil.AssertNoError(t, err)
	
	contentStr := string(content)
	testutil.AssertContains(t, contentStr, "INFO")
	testutil.AssertContains(t, contentStr, "User logged in")
	testutil.AssertContains(t, contentStr, "testuser")
	testutil.AssertContains(t, contentStr, "login")
	testutil.AssertContains(t, contentStr, "success")
}

// TestFileLoggerGetLogPath tests getting the log path
func TestFileLoggerGetLogPath(t *testing.T) {
	fl, err := NewFileLogger()
	testutil.AssertNoError(t, err)
	defer fl.Close()
	
	logPath := fl.GetLogPath()
	testutil.AssertTrue(t, logPath != "", "Log path should not be empty")
	testutil.AssertTrue(t, strings.Contains(logPath, ".devplatform"), "Log path should contain .devplatform")
	testutil.AssertTrue(t, strings.Contains(logPath, "logs"), "Log path should contain logs")
}

// TestFileLoggerRotateLogs tests log rotation
func TestFileLoggerRotateLogs(t *testing.T) {
	// Create a temporary log directory for testing
	tmpDir := t.TempDir()
	
	fl := &FileLogger{
		logDir: tmpDir,
	}
	
	// Create some old log files with valid filenames
	for i := 1; i <= 12; i++ {
		filename := filepath.Join(tmpDir, "devplatform_2024-01-01_10-00-00.log")
		if i < 10 {
			filename = filepath.Join(tmpDir, "devplatform_2024-01-0"+string(rune('0'+i))+"_10-00-00.log")
		} else {
			filename = filepath.Join(tmpDir, "devplatform_2024-01-"+string(rune('0'+i-10))+string(rune('0'))+"_10-00-00.log")
		}
		err := os.WriteFile(filename, []byte("test"), 0644)
		testutil.AssertNoError(t, err)
	}
	
	// Rotate logs
	err := fl.RotateLogs()
	testutil.AssertNoError(t, err)
	
	// Count remaining files
	entries, err := os.ReadDir(tmpDir)
	testutil.AssertNoError(t, err)
	
	count := 0
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "devplatform_") && strings.HasSuffix(entry.Name(), ".log") {
			count++
		}
	}
	
	// Should have maxLogFiles - 1 files (we delete oldest to make room for new one)
	testutil.AssertTrue(t, count <= maxLogFiles, "Should have at most maxLogFiles after rotation")
}

// TestFileLoggerClose tests closing the file logger
func TestFileLoggerClose(t *testing.T) {
	fl, err := NewFileLogger()
	testutil.AssertNoError(t, err)
	
	// Close should not error
	err = fl.Close()
	testutil.AssertNoError(t, err)
	
	// Second close might error on some systems, so we just check it doesn't panic
	_ = fl.Close()
}

// TestFileLoggerWriteAfterClose tests writing after close
func TestFileLoggerWriteAfterClose(t *testing.T) {
	fl, err := NewFileLogger()
	testutil.AssertNoError(t, err)
	
	// Close the logger
	fl.Close()
	
	// Try to write - should error
	_, err = fl.Write([]byte("test"))
	testutil.AssertError(t, err)
	// Error message varies by OS, just check it's an error
}

// TestFileLoggerWriteJSONAfterClose tests writing JSON after close
func TestFileLoggerWriteJSONAfterClose(t *testing.T) {
	fl, err := NewFileLogger()
	testutil.AssertNoError(t, err)
	
	// Close the logger
	fl.Close()
	
	// Try to write JSON - should error
	err = fl.WriteJSON("INFO", "test", map[string]interface{}{})
	testutil.AssertError(t, err)
	// Error message varies by OS, just check it's an error
}

// TestFileLoggerRotateLogsEmptyDir tests rotating logs in empty directory
func TestFileLoggerRotateLogsEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	
	fl := &FileLogger{
		logDir: tmpDir,
	}
	
	// Rotate logs in empty directory - should not error
	err := fl.RotateLogs()
	testutil.AssertNoError(t, err)
}

// TestFileLoggerRotateLogsNonexistentDir tests rotating logs when directory doesn't exist
func TestFileLoggerRotateLogsNonexistentDir(t *testing.T) {
	fl := &FileLogger{
		logDir: "/nonexistent/directory/that/does/not/exist",
	}
	
	// Rotate logs - should not error (directory doesn't exist yet)
	err := fl.RotateLogs()
	testutil.AssertNoError(t, err)
}

// TestFileLoggerGetCurrentLogFileBeforeOpen tests getting current log file before opening
func TestFileLoggerGetCurrentLogFileBeforeOpen(t *testing.T) {
	fl := &FileLogger{}
	
	logFile := fl.GetCurrentLogFile()
	testutil.AssertEqual(t, "", logFile)
}
