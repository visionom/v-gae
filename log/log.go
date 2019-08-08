package log

import (
	"sync"
)

type Logger interface {
	D(args ...interface{})
	Df(format string, args ...interface{})
	E(args ...interface{})
	Ef(format string, args ...interface{})
	I(args ...interface{})
	If(format string, args ...interface{})
	W(args ...interface{})
	Wf(format string, args ...interface{})
	V(args ...interface{})
	Vf(format string, args ...interface{})
	Flush()
}

var logger Logger

var instance Logger
var once sync.Once

func GetLogger() Logger {
	if instance == nil {
		once.Do(func() {
			instance = NewBaseLogger()
		})
	}
	return instance
}

func D(args ...interface{}) {
	GetLogger().D(args...)
}

func Df(format string, args ...interface{}) {
	GetLogger().Df(format, args...)
}

func E(args ...interface{}) {
	GetLogger().E(args...)
}

func Ef(format string, args ...interface{}) {
	GetLogger().Ef(format, args...)
}

func I(args ...interface{}) {
	GetLogger().I(args...)
}

func If(format string, args ...interface{}) {
	GetLogger().If(format, args...)
}

func W(args ...interface{}) {
	GetLogger().W(args...)
}

func Wf(format string, args ...interface{}) {
	GetLogger().Wf(format, args...)
}

func V(args ...interface{}) {
	GetLogger().V(args...)
}

func Vf(format string, args ...interface{}) {
	GetLogger().Vf(format, args...)
}

func Flush() {
	GetLogger().Flush()
}
