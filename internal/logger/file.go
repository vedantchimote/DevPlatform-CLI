package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	maxLogFiles = 10
	logDirName  = ".devplatform"
	logSubDir   = "logs"
)

// FileLogger handles file-based logging with rotation
type FileLogger struct {
	logDir  string
	logFile *os.File
}

// NewFileLogger creates a new file logger
func NewFileLogger() (*FileLogger, error) {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	
	// Create log directory path
	logDir := filepath.Join(homeDir, logDirName, logSubDir)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	
	fl := &FileLogger{
		logDir: logDir,
	}
	
	// Rotate logs before opening new file
	if err := fl.RotateLogs(); err != nil {
		return nil, fmt.Errorf("failed to rotate logs: %w", err)
	}
	
	// Open new log file
	if err := fl.openLogFile(); err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	
	return fl, nil
}

// openLogFile opens a new log file with timestamp
func (fl *FileLogger) openLogFile() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("devplatform_%s.log", timestamp)
	filepath := filepath.Join(fl.logDir, filename)
	
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	
	fl.logFile = file
	return nil
}

// Write implements io.Writer interface
func (fl *FileLogger) Write(p []byte) (n int, err error) {
	if fl.logFile == nil {
		return 0, fmt.Errorf("log file not open")
	}
	return fl.logFile.Write(p)
}

// WriteJSON writes a structured log entry as JSON
func (fl *FileLogger) WriteJSON(level string, msg string, fields map[string]interface{}) error {
	if fl.logFile == nil {
		return fmt.Errorf("log file not open")
	}
	
	entry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"message":   msg,
	}
	
	// Add fields
	for k, v := range fields {
		entry[k] = v
	}
	
	// Marshal to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}
	
	// Write to file
	_, err = fl.logFile.Write(append(data, '\n'))
	return err
}

// Close closes the log file
func (fl *FileLogger) Close() error {
	if fl.logFile != nil {
		return fl.logFile.Close()
	}
	return nil
}

// GetLogPath returns the path to the log directory
func (fl *FileLogger) GetLogPath() string {
	return fl.logDir
}

// RotateLogs removes old log files, keeping only the most recent maxLogFiles
func (fl *FileLogger) RotateLogs() error {
	// Read directory
	entries, err := os.ReadDir(fl.logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist yet, nothing to rotate
		}
		return fmt.Errorf("failed to read log directory: %w", err)
	}
	
	// Filter log files
	var logFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "devplatform_") && strings.HasSuffix(entry.Name(), ".log") {
			logFiles = append(logFiles, entry.Name())
		}
	}
	
	// If we have more than maxLogFiles, delete the oldest ones
	if len(logFiles) >= maxLogFiles {
		// Sort by name (which includes timestamp)
		sort.Strings(logFiles)
		
		// Delete oldest files
		numToDelete := len(logFiles) - maxLogFiles + 1
		for i := 0; i < numToDelete; i++ {
			filepath := filepath.Join(fl.logDir, logFiles[i])
			if err := os.Remove(filepath); err != nil {
				// Log error but continue
				fmt.Fprintf(os.Stderr, "Warning: failed to delete old log file %s: %v\n", filepath, err)
			}
		}
	}
	
	return nil
}

// GetCurrentLogFile returns the path to the current log file
func (fl *FileLogger) GetCurrentLogFile() string {
	if fl.logFile == nil {
		return ""
	}
	return fl.logFile.Name()
}
