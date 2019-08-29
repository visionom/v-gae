package ifs

import (
	"github.com/kataras/iris/context"
)

type WebServer interface {
	AddHandlers(path string, handlers map[string]context.Handler) error
	Run() error
}
