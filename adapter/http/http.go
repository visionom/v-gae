package http

import (
	"github.com/visionom/v-gae/adapter/http/config"
	"github.com/visionom/v-gae/adapter/http/interfaces"
	"github.com/visionom/v-gae/adapter/http/usecase"
)

func New(config *config.HTTPConfig) interfaces.WebServer {
	return usecase.NewWebServer(config)
}
