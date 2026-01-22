package logger

import (
	"log"
	"os"
	"sync"
)

// LogLevel represents the log level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warn    *log.Logger
	err     *log.Logger
	fatal   *log.Logger
	level   LogLevel
	errFile *os.File
	mu      sync.Mutex
}

var instance *Logger

func Init() *Logger {
	if instance != nil {
		return instance
	}

	// Create logs directory if not exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// Open error log file
	errFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	instance = &Logger{
		debug:   log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		info:    log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		warn:    log.New(os.Stdout, "[WARN] ", log.LstdFlags),
		err:     log.New(errFile, "[ERROR] ", log.LstdFlags),
		fatal:   log.New(errFile, "[FATAL] ", log.LstdFlags),
		level:   INFO, // Default to INFO level
		errFile: errFile,
	}

	return instance
}

func GetLogger() *Logger {
	if instance == nil {
		return Init()
	}
	return instance
}

// SetLogLevel sets the logging level
func (l *Logger) SetLogLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLogLevel returns the current logging level
func (l *Logger) GetLogLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= DEBUG {
		l.debug.Printf(msg+"\n", args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= INFO {
		l.info.Printf(msg+"\n", args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= WARN {
		l.warn.Printf(msg+"\n", args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= ERROR {
		l.err.Printf(msg+"\n", args...)
	}
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fatal.Printf(msg+"\n", args...)
	os.Exit(1)
}

func (l *Logger) Close() error {
	return l.errFile.Close()
}
