package logger

import (
	"log"
	"os"
)

const (
	FATAL = iota
	ERROR
	INFO
)

// Config represents a logger configuration
//
// The defined levels are FATAL, ERROR, WARN, INFO, DEBUG  with values 0-4 in the same order. Set any other value to turn off logging.
//
// Level represents the maximum log level. For example, if it is set to 2, all logs below level 3 will be printed i.e. FATAL, ERROR, and WARN
//
// ** Please Note: The FATAL level causes the program to panic **
type Config struct {
	Name string
}

// Logger methods to call on a logger instance
type Logger interface {
	Info(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	print(caller int, msgs ...interface{})
}

// NewLogger retuns a logger instance
func NewLogger(opts Config) Logger {
	return &Config{opts.Name}
}

func (l *Config) print(caller int, msg ...interface{}) {
	name := l.Name

	var levelName string
	switch caller {
	case FATAL:
		levelName = "[FATAL]"
	case ERROR:
		levelName = "[ERROR]"
	case INFO:
		levelName = "[INFO]"
	default:
		return
	}

	var prefix []interface{}
	prefix = append(prefix, levelName, name)

	msg = append(prefix, msg...)
	log.Println(msg...)
}

// Info log an info statement
func (l *Config) Info(msg ...interface{}) {
	l.print(INFO, msg...)
}

// Error log a error statement
func (l *Config) Error(msg ...interface{}) {
	l.print(ERROR, msg...)
}

// Fatal log a fatal statement causing program to panic
func (l *Config) Fatal(msg ...interface{}) {
	l.print(FATAL, msg...)
	os.Exit(1)
}

// Child create a child logger
//
// Defaults to the values set for parent logger
func (l *Config) Child(opts Config) Logger {
	name := "(" + opts.Name + ")"
	if l.Name != "" {
		name = "<" + l.Name + ">" + name
	}
	childLogger := NewLogger(Config{name})
	return childLogger
}
