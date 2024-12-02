package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/porseOnline/config"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

type LogOutput interface {
	Write(entry LogEntry) error
	Close() error
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	Fields    Fields    `json:"fields,omitempty"`
}

type Fields map[string]interface{}

type Logger struct {
	output   LogOutput
	minLevel LogLevel
	mu       sync.Mutex
}

type ConsoleOutput struct {
	writer io.Writer
}

type FileOutput struct {
	file *os.File
}

type ElasticsearchOutput struct {
	url      string
	index    string
	username string
	password string
}

var defaultLogger *Logger
var once sync.Once

func InitLogger(c config.Config) error {
	var err error
	once.Do(func() {
		_logLevel, err := parseLogLevel(c.Logger.Level)
		if err != nil {
			log.Fatal("wrong log level")
		}

		var output LogOutput
		switch c.Logger.Output {
		case "console":
			output = NewConsoleOutput(os.Stdout)
		case "file":
			output, err = NewFileOutput(c.Logger.Path)
		case "elasticsearch":
			output = NewElasticsearchOutput(c.Elasticsearch.Host, c.Elasticsearch.Index,
				c.Elasticsearch.Username, c.Elasticsearch.Password)
		}
		if err != nil {
			return
		}
		if output != nil {
			defaultLogger = &Logger{
				minLevel: _logLevel,
				output:   output,
			}

		}

	})
	return err
}

func NewConsoleOutput(w io.Writer) LogOutput {
	return &ConsoleOutput{writer: w}
}

func (c *ConsoleOutput) Write(entry LogEntry) error {
	output := fmt.Sprintf("[%s] %s: %s\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level.String(),
		entry.Message)
	_, err := fmt.Fprint(c.writer, output)
	return err
}

func (c *ConsoleOutput) Close() error {
	return nil
}

func NewFileOutput(path string) (LogOutput, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileOutput{file: file}, nil
}

func (f *FileOutput) Write(entry LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = f.file.Write(append(data, '\n'))
	return err
}

func (f *FileOutput) Close() error {
	return f.file.Close()
}

func NewElasticsearchOutput(url, index, username, password string) LogOutput {
	return &ElasticsearchOutput{
		url:      url,
		index:    index,
		username: username,
		password: password,
	}
}

func (e *ElasticsearchOutput) Write(entry LogEntry) error {
	// Implement ES client logic here
	// Use elastic client to send log entry to ES
	return nil
}

func (e *ElasticsearchOutput) Close() error {
	// Implement ES client cleanup
	return nil
}

func Debug(msg string, fields Fields) {
	defaultLogger.log(DEBUG, msg, fields)
}

func Info(msg string, fields Fields) {
	defaultLogger.log(INFO, msg, fields)
}

func Warn(msg string, fields Fields) {
	defaultLogger.log(WARN, msg, fields)
}

func Error(msg string, fields Fields) {
	defaultLogger.log(ERROR, msg, fields)
}

func Fatal(msg string, fields Fields) {
	defaultLogger.log(FATAL, msg, fields)
	os.Exit(1)
}

func (l *Logger) log(level LogLevel, msg string, fields Fields) {
	if level < l.minLevel {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.output.Write(entry)
}

func parseLogLevel(level string) (LogLevel, error) {
	switch level {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		return DEBUG, fmt.Errorf("invalid log level: %s", level)
	}
}
