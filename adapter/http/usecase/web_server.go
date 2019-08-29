package usecase

import (
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/logger"

	"github.com/visionom/v-gae/log"

	"github.com/visionom/v-gae/adapter/http/config"
	"github.com/visionom/v-gae/adapter/http/domain"
	"github.com/visionom/v-gae/adapter/http/ifs"
)

type IrisWebServer struct {
	bind     string
	hostname string
	app      *iris.Application
	root     iris.Party
}

func NewWebServer(config *config.HTTPConfig) ifs.WebServer {

	_logger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		LogFunc: func(now time.Time,
			latency time.Duration,
			status,
			ip,
			method,
			path string,
			message interface{},
			headerMessage interface{}) {
			log.If("%s %s %s %s %s", status, latency, ip, method, path)
			log.Vf("message %v", message)
			log.Vf("header  %v", headerMessage)
		},
	})
	app := iris.New()
	app.Use(_logger)
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	var root iris.Party
	if config.Root != "" {
		root = app.Party(config.Root)
	}
	root = app
	return &IrisWebServer{
		bind: config.Bind,
		app:  app,
		root: root,
	}
}

func (s *IrisWebServer) AddHandlers(path string, handlers map[string]context.Handler) error {
	p := s.root.Party(path)
	{
		for handlerPath, handler := range handlers {
			p.Any(handlerPath, handler)
		}
	}
	return nil
}

func (s *IrisWebServer) Run() error {
	return s.app.Run(iris.Addr(s.bind))
}

func JsonRes(ctx iris.Context, ticketID string, err error,
	data interface{}) {
	ctx.StatusCode(iris.StatusOK)
	code := iris.StatusOK
	msg := "success"
	if err != nil {
		code = iris.StatusInternalServerError
		msg = err.Error()
	}

	res := domain.NewRsp(ticketID, code, msg, data)
	ctx.JSON(res)
}

func ParseJsonErr(ctx iris.Context, ticketID string, err error) {
	code := iris.StatusBadRequest
	msg := "parse request json Error"
	res := domain.NewRsp(ticketID, code, msg, err)
	ctx.JSON(res)
}
