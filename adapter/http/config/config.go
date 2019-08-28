package config

//	"fmt"

//disable local config file
type HTTPConfig struct {
	Bind string `json:"bind"`
	Root string `json:"root"`
}

func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Bind: "127.0.0.1:8080",
		Root: "/api",
	}
}

func (conf *HTTPConfig) Valid() bool {
	return true
}
