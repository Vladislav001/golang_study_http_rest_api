package main

// @see https://www.youtube.com/watch?v=LxJLuW5aUDQ&list=PLehOyJfJkFkJ5m37b4oWh783yzVlHdnUH

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Vladislav001/golang_study_http_rest_api/internal/app/apiserver"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
