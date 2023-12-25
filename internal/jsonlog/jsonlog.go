package jsonlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
		mu:       sync.Mutex{},
	}
}

func (l *Logger) LogInfo(message string, properties map[string]string) {
	l.log(LevelInfo, message, properties)
}

func (l *Logger) LogError(err error, properties map[string]string) {
	l.log(LevelError, err.Error(), properties)
}

func (l *Logger) LogFatal(err error, properties map[string]string) {
	l.log(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) Write(message []byte) (int, error) {
	return l.log(LevelError, string(message), nil)
}

func (l *Logger) log(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(fmt.Sprintf("%s: unable to marshal log message: %s", LevelError.String(), err))
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}
