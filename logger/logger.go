package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

func (l LogLevel) String() string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level   LogLevel
	verbose bool
	mu      sync.Mutex
	output  io.Writer
}

var global *Logger

func New(level LogLevel, verbose bool) *Logger {
	return &Logger{
		level:   level,
		verbose: verbose,
		output:  os.Stdout,
	}
}

func SetGlobal(l *Logger) {
	global = l
}

func GetGlobal() *Logger {
	if global == nil {
		global = New(LevelInfo, false)
	}
	return global
}

func (l *Logger) GetLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *Logger) IsVerbose() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.verbose
}

func SetLevel(level LogLevel) {
	if global != nil {
		global.mu.Lock()
		global.level = level
		global.mu.Unlock()
	}
}

func SetVerbose(verbose bool) {
	if global != nil {
		global.mu.Lock()
		global.verbose = verbose
		global.mu.Unlock()
	}
}

func SetOutput(w io.Writer) {
	if global != nil {
		global.mu.Lock()
		global.output = w
		global.mu.Unlock()
	}
}

func ParseLevel(levelStr string) (LogLevel, error) {
	switch levelStr {
	case "error":
		return LevelError, nil
	case "warn":
		return LevelWarn, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	default:
		return LevelInfo, fmt.Errorf("invalid log level: %s", levelStr)
	}
}

func SetLevelString(levelStr string) error {
	level, err := ParseLevel(levelStr)
	if err != nil {
		return err
	}
	SetLevel(level)
	return nil
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level > l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := ""
	if l.verbose {
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z ")[:20] + " "
	}

	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "%s[%s] %s\n", timestamp, level.String(), message)
}

func Debug(format string, args ...interface{}) {
	GetGlobal().log(LevelDebug, format, args...)
}

func Info(format string, args ...interface{}) {
	GetGlobal().log(LevelInfo, format, args...)
}

func Warn(format string, args ...interface{}) {
	GetGlobal().log(LevelWarn, format, args...)
}

func Error(format string, args ...interface{}) {
	GetGlobal().log(LevelError, format, args...)
}
