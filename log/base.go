package log

import (
	"fmt"
	"log"
)

const (
	DEBUG   = "DBUG"
	ERROR   = "ERRO"
	INFO    = "INFO"
	WARNING = "WARN"
	VERBOSE = "VERB"
)

type BaseLogger struct {
	debug bool
	erro  bool
	info  bool
	warn  bool
	ver   bool
}

func NewBaseLogger() Logger {
	return &BaseLogger{
		debug: true,
		erro:  true,
		info:  true,
		warn:  true,
		ver:   true,
	}
}

func (l *BaseLogger) printf(level, format string, v ...interface{}) {
	log.Printf("[%s] "+format, append([]interface{}{level}, v...)...)
}

func (l *BaseLogger) print(level string, v ...interface{}) {
	pre := fmt.Sprintf("[%s]", level)
	log.Print(append([]interface{}{pre}, v...)...)
}

func (l *BaseLogger) D(args ...interface{}) {
	if l.debug {
		l.print(DEBUG, args...)
	}
}

func (l *BaseLogger) Df(format string, args ...interface{}) {
	if l.debug {
		l.printf(DEBUG, format, args...)
	}
}

func (l *BaseLogger) E(args ...interface{}) {
	if l.erro {
		l.print(ERROR, args...)
	}
}

func (l *BaseLogger) Ef(format string, args ...interface{}) {
	if l.erro {
		l.printf(ERROR, format, args...)
	}
}

func (l *BaseLogger) I(args ...interface{}) {
	if l.info {
		l.print(INFO, args...)
	}
}

func (l *BaseLogger) If(format string, args ...interface{}) {
	if l.info {
		l.printf(INFO, format, args...)
	}
}

func (l *BaseLogger) W(args ...interface{}) {
	if l.warn {
		l.print(WARNING, args...)
	}
}

func (l *BaseLogger) Wf(format string, args ...interface{}) {
	if l.warn {
		l.printf(WARNING, format, args...)
	}
}

func (l *BaseLogger) V(args ...interface{}) {
	if l.ver {
		l.print(VERBOSE, args...)
	}
}

func (l *BaseLogger) Vf(format string, args ...interface{}) {
	if l.ver {
		l.printf(VERBOSE, format, args...)
	}
}

func (l *BaseLogger) Flush() {
}
