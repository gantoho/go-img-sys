package logger

import (
	"io"
	"log"
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Config struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      LogLevel
}

var DefaultConfig = Config{
	Filename:   "logs/error.log",
	MaxSize:    100,
	MaxBackups: 7,
	MaxAge:     30,
	Compress:   true,
	Level:      INFO,
}

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warn    *log.Logger
	err     *log.Logger
	fatal   *log.Logger
	level   LogLevel
	closers []io.Closer
	mu      sync.Mutex
}

var instance *Logger

func Init() *Logger {
	return InitWithConfig(DefaultConfig)
}

func InitWithConfig(cfg Config) *Logger {
	if instance != nil {
		return instance
	}

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	stdoutWriter := os.Stdout

	errWriter := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	instance = &Logger{
		debug:   log.New(stdoutWriter, "[DEBUG] ", log.LstdFlags|log.Lmicroseconds),
		info:    log.New(stdoutWriter, "[INFO] ", log.LstdFlags|log.Lmicroseconds),
		warn:    log.New(io.MultiWriter(stdoutWriter, errWriter), "[WARN] ", log.LstdFlags|log.Lmicroseconds),
		err:     log.New(errWriter, "[ERROR] ", log.LstdFlags|log.Lmicroseconds),
		fatal:   log.New(errWriter, "[FATAL] ", log.LstdFlags|log.Lmicroseconds),
		level:   cfg.Level,
		closers: []io.Closer{errWriter},
	}

	return instance
}

func GetLogger() *Logger {
	if instance == nil {
		return Init()
	}
	return instance
}

func (l *Logger) SetLogLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

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
	panic(msg)
}

func (l *Logger) Close() error {
	for _, c := range l.closers {
		c.Close()
	}
	return nil
}
