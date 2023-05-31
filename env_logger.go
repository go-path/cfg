package cfg

import (
	"log"
	"os"
)

// Logger is implemented by any logging system that is used for standard logs.
type Logger interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
}

func (c *Env) SetLogger(logger Logger) {
	c.logger = logger
}

type defaultLog struct {
	*log.Logger
}

func defaultLogger() *defaultLog {
	return &defaultLog{Logger: log.New(os.Stderr, "cfg ", log.LstdFlags)}
}

func (l *defaultLog) Error(f string, v ...interface{}) {
	l.Printf("ERROR: "+f, v...)
}

func (l *defaultLog) Info(f string, v ...interface{}) {
	l.Printf("INFO: "+f, v...)
}

func (l *defaultLog) Warn(f string, v ...interface{}) {
	l.Printf("WARN: "+f, v...)
}
