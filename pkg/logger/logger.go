// Package logger provides a flexible logging system with multiple outputs and log levels
package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// LogLevel represents the severity of the log message
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
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// LogWriter interface for different output destinations
type LogWriter interface {
	Write(entry *LogEntry) error
	Close() error
}

// LogEntry represents a single log message
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Logger is the main logger struct
type Logger struct {
	writers []LogWriter
	level   LogLevel
	mutex   sync.Mutex
}

// FileWriter implements LogWriter for file output
type FileWriter struct {
	file *os.File
}

// StdoutWriter implements LogWriter for console output
type StdoutWriter struct{}

// ElasticsearchWriter implements LogWriter for Elasticsearch output
type ElasticsearchWriter struct {
	url      string
	index    string
	username string
	password string
}

// NewLogger creates a new Logger instance
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:   level,
		writers: make([]LogWriter, 0),
	}
}

// AddWriter adds a new output destination
func (l *Logger) AddWriter(writer LogWriter) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.writers = append(l.writers, writer)
}

// log handles the actual logging
func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	for _, writer := range l.writers {
		if err := writer.Write(entry); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log entry: %v\n", err)
		}
	}

	if level == FATAL {
		for _, writer := range l.writers {
			writer.Close()
		}
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.log(DEBUG, message, fields)
}

// Info logs an info message
func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log(INFO, message, fields)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log(WARN, message, fields)
}

// Error logs an error message
func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log(ERROR, message, fields)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(message string, fields map[string]interface{}) {
	l.log(FATAL, message, fields)
}

// NewFileWriter creates a new file writer
func NewFileWriter(filepath string) (*FileWriter, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileWriter{file: file}, nil
}

func (fw *FileWriter) Write(entry *LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = fw.file.Write(append(data, '\n'))
	return err
}

func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

// NewStdoutWriter creates a new stdout writer
func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

func (sw *StdoutWriter) Write(entry *LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stdout, string(data))
	return err
}

func (sw *StdoutWriter) Close() error {
	return nil
}

// NewElasticsearchWriter creates a new Elasticsearch writer
func NewElasticsearchWriter(url, index, username, password string) *ElasticsearchWriter {
	return &ElasticsearchWriter{
		url:      url,
		index:    index,
		username: username,
		password: password,
	}
}

func (ew *ElasticsearchWriter) Write(entry *LogEntry) error {
	// Implement Elasticsearch writing logic here
	// This is a placeholder - you'll need to implement the actual ES client logic
	return nil
}

func (ew *ElasticsearchWriter) Close() error {
	// Implement cleanup logic if needed
	return nil
}
