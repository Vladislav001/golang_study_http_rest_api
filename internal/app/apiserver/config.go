package apiserver

import "github.com/Vladislav001/golang_study_http_rest_api/internal/app/store"

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLever string `toml:"log_level"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		LogLever: "debug",
		Store:    store.NewConfig(),
	}
}
