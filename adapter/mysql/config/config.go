package config

//	"fmt"

//disable local config file
type DBConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Name   string `json:"name"`
	User   string `json:"user"`
	Passwd string `json:"passwd,pwd"`
}

func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		Host:   "127.0.0.1",
		Port:   3306,
		Name:   "test",
		User:   "root",
		Passwd: "",
	}
}

func (conf *DBConfig) Valid() bool {
	return true
}
