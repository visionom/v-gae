package http

import (
	"github.com/visionom/v-gae/adapter/http/config"
	"github.com/visionom/v-gae/adapter/http/ifs"
	"github.com/visionom/v-gae/adapter/http/usecase"
)

func New(config *config.HTTPConfig) ifs.WebServer {
	return usecase.NewWebServer(config)
}
